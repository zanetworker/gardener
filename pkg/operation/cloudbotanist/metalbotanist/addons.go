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
	"github.com/gardener/gardener/pkg/operation/common"
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
