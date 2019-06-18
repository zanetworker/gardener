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

package metalbotanist

import (
	"encoding/base64"
	"fmt"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	"github.com/gardener/gardener/pkg/operation/common"
)

// FIXME: Check if we need this function
// GenerateCloudProviderConfig generates the Metal cloud provider config.
// See this for more details:
// https://github.com/metalhost/metal-ccm
func (b *MetalBotanist) GenerateCloudProviderConfig() (string, error) {
	return "", nil
}

// RefreshCloudProviderConfig refreshes the cloud provider credentials in the existing cloud
// provider config.
// Not needed on Metal (cloud provider config does not contain the credentials), hence, the
// original is returned back.
func (b *MetalBotanist) RefreshCloudProviderConfig(currentConfig map[string]string) map[string]string {
	return currentConfig
}

// FIXME: Check if this function is adequate
// GenerateKubeAPIServerServiceConfig generates the cloud provider specific values which are required to render the
// Service manifest of the kube-apiserver-service properly.
func (b *MetalBotanist) GenerateKubeAPIServerServiceConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"enableCSI": true,
	}, nil
}

// GenerateKubeAPIServerExposeConfig defines the cloud provider specific values which configure how the kube-apiserver
// is exposed to the public.
func (b *MetalBotanist) GenerateKubeAPIServerExposeConfig() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// GenerateKubeAPIServerConfig generates the cloud provider specific values which are required to render the
// Deployment manifest of the kube-apiserver properly.
func (b *MetalBotanist) GenerateKubeAPIServerConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"environment": getMetalCredentialsEnvironment(),
	}, nil
}

// GenerateCloudControllerManagerConfig generates the cloud provider specific values which are required to
// render the Deployment manifest of the cloud-controller-manager properly.
func (b *MetalBotanist) GenerateCloudControllerManagerConfig() (map[string]interface{}, string, error) {
	chartName := "metal-cloud-controller-manager"
	conf := map[string]interface{}{}
	newConf, err := b.InjectSeedShootImages(conf, common.MetalControllerManagerImageName)
	if err != nil {
		return conf, chartName, err
	}

	return newConf, chartName, nil
}

// FIXME: Check if this function is adequate
// GenerateKubeControllerManagerConfig generates the cloud provider specific values which are required to
// render the Deployment manifest of the kube-controller-manager properly.
func (b *MetalBotanist) GenerateKubeControllerManagerConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"enableCSI": true,
	}, nil
}

// GenerateKubeSchedulerConfig generates the cloud provider specific values which are required to render the
// Deployment manifest of the kube-scheduler properly.
func (b *MetalBotanist) GenerateKubeSchedulerConfig() (map[string]interface{}, error) {
	return nil, nil
}

// FIXME: Check if this function is adequate
// GenerateCSIConfig generates the configuration for CSI charts
func (b *MetalBotanist) GenerateCSIConfig() (map[string]interface{}, error) {
	return nil, nil
}

// maps are mutable, so it's safer to create a new instance
func getMetalCredentialsEnvironment() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name": "METAL_API_KEY",
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"key":  MetalAPIKey,
					"name": gardencorev1alpha1.SecretNameCloudProvider,
				},
			},
		},
	}
}

// GenerateEtcdBackupConfig returns the etcd backup configuration for the etcd Helm chart.
func (b *MetalBotanist) GenerateEtcdBackupConfig() (map[string][]byte, map[string]interface{}, error) {
	secretData := map[string][]byte{}
	backupConfigData := map[string]interface{}{}

	return secretData, backupConfigData, nil
}

// GenerateETCDStorageClassConfig generates values which are required to create etcd volume storageclass properly.
func (b *MetalBotanist) GenerateETCDStorageClassConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":        "gardener.cloud-fast",
		"capacity":    "25Gi",
		"provisioner": "net.metal.csi",
		"parameters": map[string]interface{}{
			"plan": "standard",
		},
	}
}

// DeployCloudSpecificControlPlane does any last minute updates
func (b *MetalBotanist) DeployCloudSpecificControlPlane() error {
	return nil
}
