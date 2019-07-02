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
	"fmt"
	"github.com/gardener/gardener/pkg/operation/common"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	metalgo "github.com/metal-pod/metal-go"
	"strings"
)

// DeployKube2IAMResources - Not needed on Metal
func (b *MetalBotanist) DeployKube2IAMResources() error {
	return nil
}

// DestroyKube2IAMResources - Not needed on Metal.
func (b *MetalBotanist) DestroyKube2IAMResources() error {
	return nil
}

// GenerateKube2IAMConfig - Not needed on Metal.
func (b *MetalBotanist) GenerateKube2IAMConfig() (map[string]interface{}, error) {
	return common.GenerateAddonConfig(nil, false), nil
}

// GenerateStorageClassesConfig generates values which are required to render the chart storage-classes properly.
func (b *MetalBotanist) GenerateStorageClassesConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"StorageClasses": []map[string]interface{}{
			{
				"Name":           "csi-metal-standard",
				"IsDefaultClass": true,
				"Provisioner":    "net.metal.csi",
				"Parameters": map[string]interface{}{
					"type": "standard",
				},
			},
			{
				"Name":           "csi-metal-performance",
				"IsDefaultClass": false,
				"Provisioner":    "net.metal.csi",
				"Parameters": map[string]interface{}{
					"plan": "performance",
				},
			},
		},
	}, nil
}

// GenerateNginxIngressConfig generates values which are required to render the chart nginx-ingress properly.
func (b *MetalBotanist) GenerateNginxIngressConfig() (map[string]interface{}, error) {
	return common.GenerateAddonConfig(nil, b.Shoot.NginxIngressEnabled()), nil
}

// GenerateMetalLBConfig generates values which are required to render the chart metallb properly.
func (b *MetalBotanist) GenerateMetalLBConfig() (map[string]interface{}, error) {
	if !b.Shoot.MetalLBEnabled() {
		return common.GenerateAddonConfig(nil, false), nil
	}

	svc, err := b.createSVC()
	if err != nil {
		return nil, err
	}

	// find corresponding network
	networks := b.Shoot.Info.Spec.Addons.MetalLB.Networks
	projectId := b.Shoot.Info.Status.TechnicalID

	mlb := make(map[string]interface{})

	for _, nw := range networks {
		findReq := &metalgo.NetworkFindRequest{
			ProjectID: &projectId,
			Name:  &nw.Name,
		}
		resp, err := svc.NetworkFind(findReq)
		if err != nil {
			return nil, err
		}
		if len(resp.Networks) == 0 {
			return nil, fmt.Errorf("cannot find network with prefix %q", networks)
		}
		if len(resp.Networks) > 1 {
			return nil, fmt.Errorf("%d networks found with prefix %q", len(resp.Networks), networks)
		}
		network := resp.Networks[0]

		// acquire network IPs
		var ips []string
		for i := 0; i < nw.Count; i++ {
			req := &metalgo.IPAcquireRequest{
				Projectid: projectId,
				Networkid: *network.ID,
				Name: fmt.Sprintf("metallb-%s-%d", nw.Name, i+1),
				Description: b.Shoot.Info.Namespace,
			}
			ipa, err := svc.IPAcquire(req)
			if err != nil {
				return nil, err
			}
			ips = append(ips, *ipa.IP.Ipaddress)
		}

		// extend metallb configuration
		mlb[nw.Name] = ips
	}

	return common.GenerateAddonConfig(map[string]interface{}{
		"metallb": mlb,
	}, true), nil
}

// GenerateVPNShootConfig generate cloud-specific vpn override - Metal determines the config dynamically by querying
func (b *MetalBotanist) GenerateVPNShootConfig() (map[string]interface{}, error) {
	kv := map[string]interface{}{
		"image": "metalhost/metabot:latest", "args": []string{"ip", "4", "private", "parent", "network"},
	}
	config := map[string]interface{}{
		"initContainers": []map[string]interface{}{kv},
	}
	return config, nil
}

// Helper function to create SVC
func (b *MetalBotanist) createSVC() (*metalgo.Driver, error) {
	token := strings.TrimSpace(string(b.Shoot.Secret.Data[v1alpha1.MetalAPIKey]))
	hmac := strings.TrimSpace(string(b.Shoot.Secret.Data[v1alpha1.MetalAPIHMac]))

	u, ok := b.Shoot.Secret.Data[v1alpha1.MetalAPIURL]
	if !ok {
		return nil, fmt.Errorf("missing %s in secret", v1alpha1.MetalAPIURL)
	}
	url := strings.TrimSpace(string(u))

	return metalgo.NewDriver(url, token, hmac)
}
