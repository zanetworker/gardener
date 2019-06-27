// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hybridbotanist

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	controllermanagerfeatures "github.com/gardener/gardener/pkg/controllermanager/features"
	"github.com/gardener/gardener/pkg/features"
	"github.com/gardener/gardener/pkg/operation/common"
	"github.com/gardener/gardener/pkg/utils"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/secrets"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	audit_internal "k8s.io/apiserver/pkg/apis/audit"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
	auditv1alpha1 "k8s.io/apiserver/pkg/apis/audit/v1alpha1"
	auditv1beta1 "k8s.io/apiserver/pkg/apis/audit/v1beta1"
	auditvalidation "k8s.io/apiserver/pkg/apis/audit/validation"
)

const (
	auditPolicyConfigMapDataKey = "policy"
)

var (
	chartPathControlPlane = filepath.Join(common.ChartPath, "seed-controlplane", "charts")
	runtimeScheme         = runtime.NewScheme()
	codecs                = serializer.NewCodecFactory(runtimeScheme)
	decoder               = codecs.UniversalDecoder()
)

func init() {
	_ = auditv1alpha1.AddToScheme(runtimeScheme)
	_ = auditv1beta1.AddToScheme(runtimeScheme)
	_ = auditv1.AddToScheme(runtimeScheme)
	_ = audit_internal.AddToScheme(runtimeScheme)
}

// getResourcesForAPIServer returns the cpu and memory requirements for API server based on nodeCount
func getResourcesForAPIServer(nodeCount int) (string, string, string, string) {
	var (
		cpuRequest    string
		memoryRequest string
		cpuLimit      string
		memoryLimit   string
	)

	switch {
	case nodeCount <= 2:
		cpuRequest = "800m"
		memoryRequest = "800Mi"

		cpuLimit = "1000m"
		memoryLimit = "1200Mi"
	case nodeCount <= 10:
		cpuRequest = "1000m"
		memoryRequest = "1100Mi"

		cpuLimit = "1200m"
		memoryLimit = "1900Mi"
	case nodeCount <= 50:
		cpuRequest = "1200m"
		memoryRequest = "1600Mi"

		cpuLimit = "1500m"
		memoryLimit = "3900Mi"
	case nodeCount <= 100:
		cpuRequest = "2500m"
		memoryRequest = "5200Mi"

		cpuLimit = "3000m"
		memoryLimit = "5900Mi"
	default:
		cpuRequest = "3000m"
		memoryRequest = "5200Mi"

		cpuLimit = "4000m"
		memoryLimit = "7800Mi"
	}

	return cpuRequest, memoryRequest, cpuLimit, memoryLimit
}

// DeployETCDStorageClass create the high iops storageclass required for volume used by etcd pods in seed cluster.
func (b *HybridBotanist) DeployETCDStorageClass(ctx context.Context) error {
	storageClassConfig := b.SeedCloudBotanist.GenerateETCDStorageClassConfig()
	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, "etcd-storageclass"), b.Shoot.SeedNamespace, "etcd-storageclass", nil, storageClassConfig)
}

// DeployETCD deploys two etcd clusters via StatefulSets. The first etcd cluster (called 'main') is used for all the
// data the Shoot Kubernetes cluster needs to store, whereas the second etcd cluster (called 'events') is only used to
// store the events data. The objectstore is also set up to store the backups.
func (b *HybridBotanist) DeployETCD() error {
	secretData, backupConfigData, err := b.SeedCloudBotanist.GenerateEtcdBackupConfig()
	if err != nil {
		return err
	}

	// Some cloud botanists do not yet support backup and won't return secret data.
	if secretData != nil {
		if _, err := b.K8sSeedClient.CreateSecret(b.Shoot.SeedNamespace, common.BackupSecretName, corev1.SecretTypeOpaque, secretData, true); err != nil {
			return err
		}
	}

	storageClassConfig := b.SeedCloudBotanist.GenerateETCDStorageClassConfig()
	etcdConfig := map[string]interface{}{
		"podAnnotations": map[string]interface{}{
			"checksum/secret-etcd-ca":         b.CheckSums[gardencorev1alpha1.SecretNameCAETCD],
			"checksum/secret-etcd-server-tls": b.CheckSums["etcd-server-tls"],
			"checksum/secret-etcd-client-tls": b.CheckSums["etcd-client-tls"],
		},
		"vpa": map[string]interface{}{
			"enabled": controllermanagerfeatures.FeatureGate.Enabled(features.VPA),
		},
		"storageClassName": storageClassConfig["name"].(string),
		"storageCapacity":  storageClassConfig["capacity"].(string),
	}

	// Some cloud botanists do not yet support backup and won't return backup config data.
	if backupConfigData != nil {
		etcdConfig["backup"] = backupConfigData
		etcdConfig["podAnnotations"].(map[string]interface{})["checksum/secret-etcd-backup"] = utils.HashForMap(backupConfigData)
	}

	etcd, err := b.InjectSeedShootImages(etcdConfig, common.ETCDImageName, common.ETCDBackupRestoreImageName)
	if err != nil {
		return err
	}

	for _, role := range []string{common.EtcdRoleMain, common.EtcdRoleEvents} {
		etcd["role"] = role
		if role == common.EtcdRoleEvents {
			etcd["backup"] = map[string]interface{}{
				"storageProvider": "", // No storage provider means no backup
			}
			etcd["storageCapacity"] = b.Seed.GetValidVolumeSize("10Gi")
		} else {
			// etcd-main emits extensive (histogram) metrics
			etcd["metrics"] = "extensive"
		}

		if b.Shoot.IsHibernated {
			statefulset := &appsv1.StatefulSet{}
			if err := b.K8sSeedClient.Client().Get(context.TODO(), kutil.Key(b.Shoot.SeedNamespace, fmt.Sprintf("etcd-%s", role)), statefulset); err != nil && !apierrors.IsNotFound(err) {
				return err
			}

			if statefulset.Spec.Replicas == nil {
				etcd["replicas"] = 0
			} else {
				etcd["replicas"] = *statefulset.Spec.Replicas
			}
		}

		if err := b.ApplyChartSeed(filepath.Join(chartPathControlPlane, "etcd"), b.Shoot.SeedNamespace, fmt.Sprintf("etcd-%s", role), nil, etcd); err != nil {
			if role == common.EtcdRoleMain {
				// Since we have to update volumeClaimTemplate in existing statefulset, which is forbidden
				// by k8s. So, we have to explicitly delete the old statefulset and create new one.
				// TODO: This is backward compatibility code and should be removed in further releases.
				if apierrors.IsInvalid(err) {
					if err := b.K8sSeedClient.DeleteStatefulSet(b.Shoot.SeedNamespace, fmt.Sprintf("etcd-%s", role)); err != nil && !apierrors.IsNotFound(err) {
						return err
					}

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
					defer cancel()
					if err := b.Botanist.WaitUntilEtcdStatefulsetDeleted(ctx, role); err != nil {
						return err
					}

					if err := b.ApplyChartSeed(filepath.Join(chartPathControlPlane, "etcd"), b.Shoot.SeedNamespace, fmt.Sprintf("etcd-%s", role), nil, etcd); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
		}
		if err := b.K8sSeedClient.DeleteService(b.Shoot.SeedNamespace, fmt.Sprintf("etcd-%s", role)); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (b *HybridBotanist) deployNetworkPolicies(ctx context.Context, denyAll bool) error {
	var (
		globalNetworkPoliciesValues = map[string]interface{}{
			"blockedAddresses": b.Seed.Info.Spec.BlockCIDRs,
			"denyAll":          denyAll,
		}
		excludeNets = []gardencorev1alpha1.CIDR{}

		values            = map[string]interface{}{}
		shootCIDRNetworks = []gardencorev1alpha1.CIDR{}
	)

	for _, addr := range b.Seed.Info.Spec.BlockCIDRs {
		excludeNets = append(excludeNets, addr)
	}

	networks, err := b.Shoot.GetK8SNetworks()
	if err != nil {
		return err
	}
	if networks != nil {
		if networks.Nodes != nil {
			shootCIDRNetworks = append(shootCIDRNetworks, *networks.Nodes)
		}
		if networks.Pods != nil {
			shootCIDRNetworks = append(shootCIDRNetworks, *networks.Pods)
		}
		if networks.Services != nil {
			shootCIDRNetworks = append(shootCIDRNetworks, *networks.Services)
		}
		shootNetworkValues, err := common.ExceptNetworks(shootCIDRNetworks, excludeNets...)
		if err != nil {
			return err
		}
		values["clusterNetworks"] = shootNetworkValues
	}

	seedNetworks := b.Seed.Info.Spec.Networks
	allCIDRNetworks := append([]gardencorev1alpha1.CIDR{seedNetworks.Nodes, seedNetworks.Pods, seedNetworks.Services}, shootCIDRNetworks...)
	allCIDRNetworks = append(allCIDRNetworks, excludeNets...)

	privateNetworks, err := common.ToExceptNetworks(common.AllPrivateNetworkBlocks(), allCIDRNetworks...)
	if err != nil {
		return err
	}
	globalNetworkPoliciesValues["privateNetworks"] = privateNetworks
	values["global-network-policies"] = globalNetworkPoliciesValues

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, "network-policies"), b.Shoot.SeedNamespace, "network-policies", values, nil)
}

// DeployLimitedNetworkPolicies creates a network policies in a Shoot cluster's namespace that
// DOES NOT deny all traffic and allow certain components to use annotations to declare their desire
// to transmit/receive traffic to/from other Pods/IP addresses.
// This is needed until migration to the new policies is complete.
func (b *HybridBotanist) DeployLimitedNetworkPolicies(ctx context.Context) error {
	return b.deployNetworkPolicies(ctx, false)
}

// DeployNetworkPolicies creates a network policies in a Shoot cluster's namespace that
// deny all traffic and allow certain components to use annotations to declare their desire
// to transmit/receive traffic to/from other Pods/IP addresses.
func (b *HybridBotanist) DeployNetworkPolicies(ctx context.Context) error {
	return b.deployNetworkPolicies(ctx, true)
}

// DeployCloudProviderConfig asks the Cloud Botanist to provide the cloud specific values for the cloud
// provider configuration. It will create a ConfigMap for it and store it in the Seed cluster.
func (b *HybridBotanist) DeployCloudProviderConfig() error {
	cloudProviderConfig, err := b.ShootCloudBotanist.GenerateCloudProviderConfig()
	if err != nil {
		return err
	}
	b.CheckSums[common.CloudProviderConfigName] = computeCloudProviderConfigChecksum(cloudProviderConfig)

	defaultValues := map[string]interface{}{
		"cloudProviderConfig": cloudProviderConfig,
	}

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, common.CloudProviderConfigName), b.Shoot.SeedNamespace, common.CloudProviderConfigName, nil, defaultValues)
}

// RefreshCloudProviderConfig asks the Cloud Botanist to refresh the cloud provider config in case it stores
// the cloud provider credentials. The Cloud Botanist is expected to return the complete updated cloud config.
func (b *HybridBotanist) RefreshCloudProviderConfig() error {
	currentConfig, err := b.K8sSeedClient.GetConfigMap(b.Shoot.SeedNamespace, common.CloudProviderConfigName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	newConfigData := b.ShootCloudBotanist.RefreshCloudProviderConfig(currentConfig.Data)
	b.CheckSums[common.CloudProviderConfigName] = computeCloudProviderConfigChecksum(newConfigData[common.CloudProviderConfigMapKey])

	_, err = b.K8sSeedClient.UpdateConfigMap(b.Shoot.SeedNamespace, common.CloudProviderConfigName, newConfigData)
	return err
}

// RefreshCSIControllersChecksums updates the cloud provider checksum in the kube-controller-manager pod spec template.
func (b *HybridBotanist) RefreshCSIControllersChecksums() error {
	if _, err := b.K8sSeedClient.GetDeployment(b.Shoot.SeedNamespace, common.CSIPluginController); err != nil {
		if apierrors.IsNotFound(err) {
			return b.DeployCSIControllers()
		}
		return err
	}

	type jsonPatch struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	patch := []jsonPatch{
		{
			Op:    "replace",
			Path:  "/spec/template/metadata/annotations/checksum~1cloudprovider",
			Value: b.CheckSums[gardencorev1alpha1.SecretNameCloudProvider],
		},
	}

	body, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = b.K8sSeedClient.PatchDeployment(b.Shoot.SeedNamespace, common.CSIPluginController, body)
	return err
}

func computeCloudProviderConfigChecksum(cloudProviderConfig string) string {
	return utils.ComputeSHA256Hex([]byte(strings.TrimSpace(cloudProviderConfig)))
}

// DeployKubeAPIServerService asks the Cloud Botanist to provide the cloud specific configuration values for the
// kube-apiserver service.
func (b *HybridBotanist) DeployKubeAPIServerService() error {
	var (
		name          = "kube-apiserver-service"
		defaultValues = map[string]interface{}{}
	)

	cloudSpecificValues, err := b.SeedCloudBotanist.GenerateKubeAPIServerServiceConfig()
	if err != nil {
		return err
	}

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, name), b.Shoot.SeedNamespace, name, defaultValues, cloudSpecificValues)
}

// DeployKubeAPIServer asks the Cloud Botanist to provide the cloud specific configuration values for the
// kube-apiserver deployment.
func (b *HybridBotanist) DeployKubeAPIServer() error {
	defaultValues := map[string]interface{}{
		"etcdServicePort":   2379,
		"kubernetesVersion": b.Shoot.Info.Spec.Kubernetes.Version,
		"shootNetworks": map[string]interface{}{
			"service": b.Shoot.GetServiceNetwork(),
		},
		"seedNetworks": map[string]interface{}{
			"service": b.Seed.Info.Spec.Networks.Services,
			"pod":     b.Seed.Info.Spec.Networks.Pods,
			"node":    b.Seed.Info.Spec.Networks.Nodes,
		},
		"maxReplicas":      3,
		"securePort":       443,
		"probeCredentials": utils.EncodeBase64([]byte(fmt.Sprintf("%s:%s", b.Secrets["kubecfg"].Data[secrets.DataKeyUserName], b.Secrets["kubecfg"].Data[secrets.DataKeyPassword]))),
		"podAnnotations": map[string]interface{}{
			"checksum/secret-ca":                        b.CheckSums[gardencorev1alpha1.SecretNameCACluster],
			"checksum/secret-ca-front-proxy":            b.CheckSums[gardencorev1alpha1.SecretNameCAFrontProxy],
			"checksum/secret-kube-apiserver":            b.CheckSums[common.KubeAPIServerDeploymentName],
			"checksum/secret-kube-aggregator":           b.CheckSums["kube-aggregator"],
			"checksum/secret-kube-apiserver-kubelet":    b.CheckSums["kube-apiserver-kubelet"],
			"checksum/secret-kube-apiserver-basic-auth": b.CheckSums["kube-apiserver-basic-auth"],
			"checksum/secret-vpn-seed":                  b.CheckSums["vpn-seed"],
			"checksum/secret-vpn-seed-tlsauth":          b.CheckSums["vpn-seed-tlsauth"],
			"checksum/secret-service-account-key":       b.CheckSums["service-account-key"],
			"checksum/secret-etcd-ca":                   b.CheckSums[gardencorev1alpha1.SecretNameCAETCD],
			"checksum/secret-etcd-client-tls":           b.CheckSums["etcd-client-tls"],
		},
		"vpa": map[string]interface{}{
			"enabled": controllermanagerfeatures.FeatureGate.Enabled(features.VPA),
		},
	}
	cloudSpecificExposeValues, err := b.SeedCloudBotanist.GenerateKubeAPIServerExposeConfig()
	if err != nil {
		return err
	}
	cloudSpecificValues, err := b.ShootCloudBotanist.GenerateKubeAPIServerConfig()
	if err != nil {
		return err
	}

	if b.ShootedSeed != nil {
		var (
			apiServer  = b.ShootedSeed.APIServer
			autoscaler = apiServer.Autoscaler
		)
		defaultValues["replicas"] = *apiServer.Replicas
		defaultValues["minReplicas"] = *autoscaler.MinReplicas
		defaultValues["maxReplicas"] = autoscaler.MaxReplicas
		defaultValues["apiServerResources"] = map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "7000Mi",
			},
		}
	} else {
		deployment := &appsv1.Deployment{}
		if err := b.K8sSeedClient.Client().Get(context.TODO(), kutil.Key(b.Shoot.SeedNamespace, common.KubeAPIServerDeploymentName), deployment); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		replicas := deployment.Spec.Replicas

		// As kube-apiserver HPA manages the number of replicas, we have to maintain current number of replicas
		// otherwise keep the value to default
		if replicas != nil && *replicas > 0 {
			defaultValues["replicas"] = *replicas
		}
		// If the shoot is hibernated then we want to keep the number of replicas (scale down happens later).
		if b.Shoot.IsHibernated && (replicas == nil || *replicas == 0) {
			defaultValues["replicas"] = 0
		}

		cpuRequest, memoryRequest, cpuLimit, memoryLimit := getResourcesForAPIServer(b.Shoot.GetNodeCount())
		defaultValues["apiServerResources"] = map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    cpuLimit,
				"memory": memoryLimit,
			},
			"requests": map[string]interface{}{
				"cpu":    cpuRequest,
				"memory": memoryRequest,
			},
		}
	}

	var (
		apiServerConfig  = b.Shoot.Info.Spec.Kubernetes.KubeAPIServer
		admissionPlugins = kubernetes.GetAdmissionPluginsForVersion(b.Shoot.Info.Spec.Kubernetes.Version)
	)

	// Needed due to https://github.com/kubernetes/kubernetes/pull/73102
	if !b.Shoot.UsesCSI() {
		admissionPlugins = append(admissionPlugins, gardenv1beta1.AdmissionPlugin{Name: "PersistentVolumeLabel"})
	}

	if apiServerConfig != nil {
		defaultValues["featureGates"] = apiServerConfig.FeatureGates
		defaultValues["runtimeConfig"] = apiServerConfig.RuntimeConfig

		if apiServerConfig.OIDCConfig != nil {
			defaultValues["oidcConfig"] = apiServerConfig.OIDCConfig
		}

		for _, plugin := range apiServerConfig.AdmissionPlugins {
			pluginOverwritesDefault := false

			for i, defaultPlugin := range admissionPlugins {
				if defaultPlugin.Name == plugin.Name {
					pluginOverwritesDefault = true
					admissionPlugins[i] = plugin
					break
				}
			}

			if !pluginOverwritesDefault {
				admissionPlugins = append(admissionPlugins, plugin)
			}
		}

		if apiServerConfig.AuditConfig != nil &&
			apiServerConfig.AuditConfig.AuditPolicy != nil &&
			apiServerConfig.AuditConfig.AuditPolicy.ConfigMapRef != nil {
			auditPolicy, err := b.getAuditPolicy(apiServerConfig.AuditConfig.AuditPolicy.ConfigMapRef.Name, b.Shoot.Info.Namespace)
			if err != nil {
				return fmt.Errorf("Retrieving audit policy from the ConfigMap '%v' failed with reason '%v'", apiServerConfig.AuditConfig.AuditPolicy.ConfigMapRef.Name, err)
			}
			defaultValues["auditConfig"] = map[string]interface{}{
				"auditPolicy": auditPolicy,
			}
		}
	}
	defaultValues["admissionPlugins"] = admissionPlugins

	if b.Shoot.UsesCSI() {
		var existingFeatureGates map[string]bool
		if fg, ok := defaultValues["featureGates"]; ok {
			existingFeatureGates = fg.(map[string]bool)
		}

		featureGates, err := common.InjectCSIFeatureGates(b.ShootVersion(), existingFeatureGates)
		if err != nil {
			return err
		}
		defaultValues["featureGates"] = featureGates
	} else {
		// Needed due to https://github.com/kubernetes/kubernetes/pull/73102
		defaultValues["cloudProvider"] = b.ShootCloudBotanist.GetCloudProviderName()
		defaultValues["podAnnotations"].(map[string]interface{})["checksum/secret-cloudprovider"] = b.CheckSums[gardencorev1alpha1.SecretNameCloudProvider]
		defaultValues["podAnnotations"].(map[string]interface{})["checksum/configmap-cloud-provider-config"] = b.CheckSums[common.CloudProviderConfigName]
	}

	values, err := b.InjectSeedShootImages(defaultValues,
		common.HyperkubeImageName,
		common.VPNSeedImageName,
		common.BlackboxExporterImageName,
		common.AlpineIptablesImageName,
	)
	if err != nil {
		return err
	}

	// If shoot is hibernated we don't want the HPA to interfer with our scaling decisions.
	if b.Shoot.IsHibernated {
		if err := b.K8sSeedClient.DeleteHorizontalPodAutoscaler(b.Shoot.SeedNamespace, common.KubeAPIServerDeploymentName); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}

	if err := b.ApplyChartSeed(filepath.Join(chartPathControlPlane, common.KubeAPIServerDeploymentName), b.Shoot.SeedNamespace, common.KubeAPIServerDeploymentName, values, utils.MergeMaps(cloudSpecificExposeValues, cloudSpecificValues)); err != nil {
		return err
	}

	// Delete old network policies. This code can be removed in a future version.
	for _, name := range []string{
		"kube-apiserver-deny-blacklist",
		"kube-apiserver-allow-dns",
		"kube-apiserver-allow-etcd",
		"kube-apiserver-allow-gardener-admission-controller",
	} {
		if err := b.K8sSeedClient.DeleteNetworkPolicy(b.Shoot.SeedNamespace, name); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (b *HybridBotanist) getAuditPolicy(name, namespace string) (string, error) {
	auditPolicyCm, err := b.K8sGardenClient.GetConfigMap(namespace, name)
	if err != nil {
		return "", err
	}
	auditPolicy, ok := auditPolicyCm.Data[auditPolicyConfigMapDataKey]
	if !ok {
		return "", fmt.Errorf("Missing '.data.policy' in audit policy configmap %v/%v", namespace, name)
	}
	if len(auditPolicy) == 0 {
		return "", fmt.Errorf("Empty audit policy. Provide non-empty audit policy")
	}
	auditPolicyObj, schemaVersion, err := decoder.Decode([]byte(auditPolicy), nil, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to decode the provided audit policy err=%v", err)
	}
	auditPolicyInternal, ok := auditPolicyObj.(*audit_internal.Policy)
	if !ok {
		return "", fmt.Errorf("Failure to cast to audit Policy type: %v", schemaVersion)
	}
	errList := auditvalidation.ValidatePolicy(auditPolicyInternal)
	if len(errList) != 0 {
		return "", fmt.Errorf("Provided invalid audit policy err=%v", errList)
	}
	return auditPolicy, nil
}

// DeployKubeControllerManager asks the Cloud Botanist to provide the cloud specific configuration values for the
// kube-controller-manager deployment.
func (b *HybridBotanist) DeployKubeControllerManager() error {
	defaultValues := map[string]interface{}{
		"cloudProvider":     b.ShootCloudBotanist.GetCloudProviderName(),
		"clusterName":       b.Shoot.SeedNamespace,
		"kubernetesVersion": b.Shoot.Info.Spec.Kubernetes.Version,
		"podNetwork":        b.Shoot.GetPodNetwork(),
		"serviceNetwork":    b.Shoot.GetServiceNetwork(),
		"podAnnotations": map[string]interface{}{
			"checksum/secret-ca":                             b.CheckSums[gardencorev1alpha1.SecretNameCACluster],
			"checksum/secret-kube-controller-manager":        b.CheckSums[common.KubeControllerManagerDeploymentName],
			"checksum/secret-kube-controller-manager-server": b.CheckSums[common.KubeControllerManagerServerName],
			"checksum/secret-service-account-key":            b.CheckSums["service-account-key"],
			"checksum/secret-cloudprovider":                  b.CheckSums[gardencorev1alpha1.SecretNameCloudProvider],
			"checksum/configmap-cloud-provider-config":       b.CheckSums[common.CloudProviderConfigName],
		},
		"objectCount": b.Shoot.GetNodeCount(),
		"vpa": map[string]interface{}{
			"enabled": controllermanagerfeatures.FeatureGate.Enabled(features.VPA),
		},
	}
	cloudSpecificValues, err := b.ShootCloudBotanist.GenerateKubeControllerManagerConfig()
	if err != nil {
		return err
	}

	if b.Shoot.IsHibernated {
		replicaCount, err := common.CurrentReplicaCount(b.K8sSeedClient.Client(), b.Shoot.SeedNamespace, common.KubeControllerManagerDeploymentName)
		if err != nil {
			return err
		}
		defaultValues["replicas"] = replicaCount
	}

	controllerManagerConfig := b.Shoot.Info.Spec.Kubernetes.KubeControllerManager
	if controllerManagerConfig != nil {
		defaultValues["featureGates"] = controllerManagerConfig.FeatureGates

		if controllerManagerConfig.HorizontalPodAutoscalerConfig != nil {
			defaultValues["horizontalPodAutoscaler"] = controllerManagerConfig.HorizontalPodAutoscalerConfig
		}
	}

	values, err := b.InjectSeedShootImages(defaultValues, common.HyperkubeImageName)
	if err != nil {
		return err
	}

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, common.KubeControllerManagerDeploymentName), b.Shoot.SeedNamespace, common.KubeControllerManagerDeploymentName, values, cloudSpecificValues)
}

// DeployCloudControllerManager asks the Cloud Botanist to provide the cloud specific configuration values for the
// cloud-controller-manager deployment.
func (b *HybridBotanist) DeployCloudControllerManager() error {
	defaultValues := map[string]interface{}{
		"cloudProvider":     b.ShootCloudBotanist.GetCloudProviderName(),
		"clusterName":       b.Shoot.SeedNamespace,
		"kubernetesVersion": b.Shoot.Info.Spec.Kubernetes.Version,
		"podNetwork":        b.Shoot.GetPodNetwork(),
		"replicas":          1,
		"podAnnotations": map[string]interface{}{
			"checksum/secret-cloud-controller-manager":        b.CheckSums[common.CloudControllerManagerDeploymentName],
			"checksum/secret-cloud-controller-manager-server": b.CheckSums[common.CloudControllerManagerServerName],
			"checksum/secret-cloudprovider":                   b.CheckSums[gardencorev1alpha1.SecretNameCloudProvider],
			"checksum/configmap-cloud-provider-config":        b.CheckSums[common.CloudProviderConfigName],
		},
		"vpa": map[string]interface{}{
			"enabled": controllermanagerfeatures.FeatureGate.Enabled(features.VPA),
		},
	}
	cloudSpecificValues, chartName, err := b.ShootCloudBotanist.GenerateCloudControllerManagerConfig()
	if err != nil {
		return err
	}

	if b.ShootedSeed != nil {
		defaultValues["resources"] = map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "512Mi",
			},
		}
	}

	if b.Shoot.IsHibernated {
		replicaCount, err := common.CurrentReplicaCount(b.K8sSeedClient.Client(), b.Shoot.SeedNamespace, common.CloudControllerManagerDeploymentName)
		if err != nil {
			return err
		}
		defaultValues["replicas"] = replicaCount
	}

	cloudControllerManagerConfig := b.Shoot.Info.Spec.Kubernetes.CloudControllerManager
	if cloudControllerManagerConfig != nil {
		defaultValues["featureGates"] = cloudControllerManagerConfig.FeatureGates
	}

	values, err := b.InjectSeedShootImages(defaultValues, common.HyperkubeImageName)
	if err != nil {
		return err
	}

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, chartName), b.Shoot.SeedNamespace, common.CloudControllerManagerDeploymentName, values, cloudSpecificValues)
}

// DeployKubeScheduler asks the Cloud Botanist to provide the cloud specific configuration values for the
// kube-scheduler deployment.
func (b *HybridBotanist) DeployKubeScheduler() error {
	defaultValues := map[string]interface{}{
		"replicas":          b.Shoot.GetReplicas(1),
		"kubernetesVersion": b.Shoot.Info.Spec.Kubernetes.Version,
		"podAnnotations": map[string]interface{}{
			"checksum/secret-kube-scheduler":        b.CheckSums[common.KubeSchedulerDeploymentName],
			"checksum/secret-kube-scheduler-server": b.CheckSums[common.KubeSchedulerServerName],
		},
		"vpa": map[string]interface{}{
			"enabled": controllermanagerfeatures.FeatureGate.Enabled(features.VPA),
		},
	}
	cloudValues, err := b.ShootCloudBotanist.GenerateKubeSchedulerConfig()
	if err != nil {
		return err
	}

	if b.ShootedSeed != nil {
		defaultValues["resources"] = map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "300m",
				"memory": "512Mi",
			},
		}
	}

	schedulerConfig := b.Shoot.Info.Spec.Kubernetes.KubeScheduler
	if schedulerConfig != nil {
		defaultValues["featureGates"] = schedulerConfig.FeatureGates
	}

	values, err := b.InjectSeedShootImages(defaultValues, common.HyperkubeImageName)
	if err != nil {
		return err
	}

	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, common.KubeSchedulerDeploymentName), b.Shoot.SeedNamespace, common.KubeSchedulerDeploymentName, values, cloudValues)
}

// DeployCSIControllers deploy CSI controllers into the Shoot namespace in the Seed cluster. They include CSI plugin controller service (Provider specific),
// CSI external attacher, CSI external provisioner and CSI snapshotter.
func (b *HybridBotanist) DeployCSIControllers() error {
	csiPlugin, err := b.ShootCloudBotanist.GenerateCSIConfig()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("csi-%s", b.ShootCloudBotanist.GetCloudProviderName())
	return b.ApplyChartSeed(filepath.Join(chartPathControlPlane, name), b.Shoot.SeedNamespace, name, csiPlugin, nil)
}
