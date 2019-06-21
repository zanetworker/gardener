// +build !ignore_autogenerated

/*
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	metal "github.com/metal-pod/gardener-extension-provider-metal/pkg/apis/metal"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*CloudControllerManagerConfig)(nil), (*metal.CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(a.(*CloudControllerManagerConfig), b.(*metal.CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.CloudControllerManagerConfig)(nil), (*CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(a.(*metal.CloudControllerManagerConfig), b.(*CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ControlPlaneConfig)(nil), (*metal.ControlPlaneConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(a.(*ControlPlaneConfig), b.(*metal.ControlPlaneConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.ControlPlaneConfig)(nil), (*ControlPlaneConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(a.(*metal.ControlPlaneConfig), b.(*ControlPlaneConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Firewall)(nil), (*metal.Firewall)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Firewall_To_metal_Firewall(a.(*Firewall), b.(*metal.Firewall), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.Firewall)(nil), (*Firewall)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_Firewall_To_v1alpha1_Firewall(a.(*metal.Firewall), b.(*Firewall), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FirewallStatus)(nil), (*metal.FirewallStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_FirewallStatus_To_metal_FirewallStatus(a.(*FirewallStatus), b.(*metal.FirewallStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.FirewallStatus)(nil), (*FirewallStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_FirewallStatus_To_v1alpha1_FirewallStatus(a.(*metal.FirewallStatus), b.(*FirewallStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InfrastructureConfig)(nil), (*metal.InfrastructureConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(a.(*InfrastructureConfig), b.(*metal.InfrastructureConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.InfrastructureConfig)(nil), (*InfrastructureConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(a.(*metal.InfrastructureConfig), b.(*InfrastructureConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InfrastructureStatus)(nil), (*metal.InfrastructureStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(a.(*InfrastructureStatus), b.(*metal.InfrastructureStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.InfrastructureStatus)(nil), (*InfrastructureStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(a.(*metal.InfrastructureStatus), b.(*InfrastructureStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in *CloudControllerManagerConfig, out *metal.CloudControllerManagerConfig, s conversion.Scope) error {
	out.KubernetesConfig = in.KubernetesConfig
	return nil
}

// Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig is an autogenerated conversion function.
func Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in *CloudControllerManagerConfig, out *metal.CloudControllerManagerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in, out, s)
}

func autoConvert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in *metal.CloudControllerManagerConfig, out *CloudControllerManagerConfig, s conversion.Scope) error {
	out.KubernetesConfig = in.KubernetesConfig
	return nil
}

// Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig is an autogenerated conversion function.
func Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in *metal.CloudControllerManagerConfig, out *CloudControllerManagerConfig, s conversion.Scope) error {
	return autoConvert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in, out, s)
}

func autoConvert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in *ControlPlaneConfig, out *metal.ControlPlaneConfig, s conversion.Scope) error {
	out.CloudControllerManager = (*metal.CloudControllerManagerConfig)(unsafe.Pointer(in.CloudControllerManager))
	return nil
}

// Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig is an autogenerated conversion function.
func Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in *ControlPlaneConfig, out *metal.ControlPlaneConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in, out, s)
}

func autoConvert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in *metal.ControlPlaneConfig, out *ControlPlaneConfig, s conversion.Scope) error {
	out.CloudControllerManager = (*CloudControllerManagerConfig)(unsafe.Pointer(in.CloudControllerManager))
	return nil
}

// Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig is an autogenerated conversion function.
func Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in *metal.ControlPlaneConfig, out *ControlPlaneConfig, s conversion.Scope) error {
	return autoConvert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in, out, s)
}

func autoConvert_v1alpha1_Firewall_To_metal_Firewall(in *Firewall, out *metal.Firewall, s conversion.Scope) error {
	out.Size = in.Size
	out.Image = in.Image
	out.Networks = *(*[]string)(unsafe.Pointer(&in.Networks))
	out.Partition = in.Partition
	return nil
}

// Convert_v1alpha1_Firewall_To_metal_Firewall is an autogenerated conversion function.
func Convert_v1alpha1_Firewall_To_metal_Firewall(in *Firewall, out *metal.Firewall, s conversion.Scope) error {
	return autoConvert_v1alpha1_Firewall_To_metal_Firewall(in, out, s)
}

func autoConvert_metal_Firewall_To_v1alpha1_Firewall(in *metal.Firewall, out *Firewall, s conversion.Scope) error {
	out.Size = in.Size
	out.Image = in.Image
	out.Networks = *(*[]string)(unsafe.Pointer(&in.Networks))
	out.Partition = in.Partition
	return nil
}

// Convert_metal_Firewall_To_v1alpha1_Firewall is an autogenerated conversion function.
func Convert_metal_Firewall_To_v1alpha1_Firewall(in *metal.Firewall, out *Firewall, s conversion.Scope) error {
	return autoConvert_metal_Firewall_To_v1alpha1_Firewall(in, out, s)
}

func autoConvert_v1alpha1_FirewallStatus_To_metal_FirewallStatus(in *FirewallStatus, out *metal.FirewallStatus, s conversion.Scope) error {
	out.Succeeded = in.Succeeded
	out.MachineID = in.MachineID
	return nil
}

// Convert_v1alpha1_FirewallStatus_To_metal_FirewallStatus is an autogenerated conversion function.
func Convert_v1alpha1_FirewallStatus_To_metal_FirewallStatus(in *FirewallStatus, out *metal.FirewallStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_FirewallStatus_To_metal_FirewallStatus(in, out, s)
}

func autoConvert_metal_FirewallStatus_To_v1alpha1_FirewallStatus(in *metal.FirewallStatus, out *FirewallStatus, s conversion.Scope) error {
	out.Succeeded = in.Succeeded
	out.MachineID = in.MachineID
	return nil
}

// Convert_metal_FirewallStatus_To_v1alpha1_FirewallStatus is an autogenerated conversion function.
func Convert_metal_FirewallStatus_To_v1alpha1_FirewallStatus(in *metal.FirewallStatus, out *FirewallStatus, s conversion.Scope) error {
	return autoConvert_metal_FirewallStatus_To_v1alpha1_FirewallStatus(in, out, s)
}

func autoConvert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in *InfrastructureConfig, out *metal.InfrastructureConfig, s conversion.Scope) error {
	if err := Convert_v1alpha1_Firewall_To_metal_Firewall(&in.Firewall, &out.Firewall, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig is an autogenerated conversion function.
func Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in *InfrastructureConfig, out *metal.InfrastructureConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in, out, s)
}

func autoConvert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in *metal.InfrastructureConfig, out *InfrastructureConfig, s conversion.Scope) error {
	if err := Convert_metal_Firewall_To_v1alpha1_Firewall(&in.Firewall, &out.Firewall, s); err != nil {
		return err
	}
	return nil
}

// Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig is an autogenerated conversion function.
func Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in *metal.InfrastructureConfig, out *InfrastructureConfig, s conversion.Scope) error {
	return autoConvert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in, out, s)
}

func autoConvert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in *InfrastructureStatus, out *metal.InfrastructureStatus, s conversion.Scope) error {
	if err := Convert_v1alpha1_FirewallStatus_To_metal_FirewallStatus(&in.Firewall, &out.Firewall, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus is an autogenerated conversion function.
func Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in *InfrastructureStatus, out *metal.InfrastructureStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in, out, s)
}

func autoConvert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in *metal.InfrastructureStatus, out *InfrastructureStatus, s conversion.Scope) error {
	if err := Convert_metal_FirewallStatus_To_v1alpha1_FirewallStatus(&in.Firewall, &out.Firewall, s); err != nil {
		return err
	}
	return nil
}

// Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus is an autogenerated conversion function.
func Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in *metal.InfrastructureStatus, out *InfrastructureStatus, s conversion.Scope) error {
	return autoConvert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in, out, s)
}
