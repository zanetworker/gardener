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

package migration

import (
	"fmt"

	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"

	metalv1alpha1 "github.com/metal-pod/gardener-extension-provider-metal/pkg/apis/metal/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GardenV1beta1ShootToMetalV1alpha1InfrastructureConfig converts a garden.sapcloud.io/v1beta1.Shoot to metalv1alpha1.InfrastructureConfig.
// This function is only required temporarily for migration purposes and can be removed in the future when we switched to
// core.gardener.cloud/v1alpha1.Shoot.
func GardenV1beta1ShootToMetalV1alpha1InfrastructureConfig(shoot *gardenv1beta1.Shoot) (*metalv1alpha1.InfrastructureConfig, error) {
	if shoot.Spec.Cloud.Metal == nil {
		return nil, fmt.Errorf("shoot is not of type Metal")
	}

	return &metalv1alpha1.InfrastructureConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metalv1alpha1.SchemeGroupVersion.String(),
			Kind:       infrastructureConfig,
		},
		Firewall: metalv1alpha1.Firewall{
			// FIXME Shoot Spec must have firewall defaults
			Size:      shoot.Spec.Cloud.Metal.Workers[0].MachineType,
			Image:     "firewall-1",
			Networks:  []string{"internet-nbg-w8101", "underlay-nbg-w8101"},
			Partition: shoot.Spec.Cloud.Metal.Zones[0],
		},
	}, nil
}

// kind: InfrastructureConfig
// firewall:
//   partition: nbg-w8101
//   size: c1-xlarge-x86
//   image: firewall-1
//   networks:
// 	- internet-nbg-w8101
// 	- underlay-nbg-w8101
