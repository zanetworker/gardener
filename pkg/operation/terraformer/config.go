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

package terraformer

import (
	"context"
	"errors"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	_ "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SetVariablesEnvironment sets the provided <tfvarsEnvironment> on the Terraformer object.
func (t *Terraformer) SetVariablesEnvironment(tfvarsEnvironment map[string]string) *Terraformer {
	t.variablesEnvironment = tfvarsEnvironment
	return t
}

// SetImage sets the provided <image> on the Terraformer object.
func (t *Terraformer) SetImage(image string) *Terraformer {
	t.image = image
	return t
}

type InitializerConfig struct {
	Namespace         string
	ConfigurationName string
	VariablesName     string
	StateName         string
	IsStateEmpty      bool
}

type Initializer func(config *InitializerConfig) error

func (t *Terraformer) initializerConfig() *InitializerConfig {
	return &InitializerConfig{
		Namespace:         t.namespace,
		ConfigurationName: t.configName,
		VariablesName:     t.variablesName,
		StateName:         t.stateName,
		IsStateEmpty:      t.isStateEmpty(),
	}
}

func (t *Terraformer) InitializeWith(initializer Initializer) *Terraformer {
	if err := initializer(t.initializerConfig()); err != nil {
		t.logger.Errorf("Could not create the Terraform ConfigMaps/Secrets: %s", err.Error())
		return t
	}
	t.configurationDefined = true
	return t
}

const (
	MainKey      = "main.tf"
	VariablesKey = "variables.tf"
	TFVarsKey    = "terraform.tfvars"
	StateKey     = "terraform.tfstate"
)

func createOrUpdateConfigMap(ctx context.Context, c client.Client, namespace, name string, values map[string]string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}
	err := kutil.CreateOrUpdate(ctx, c, configMap, func() error {
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		for key, value := range values {
			configMap.Data[key] = value
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func CreateOrUpdateConfigurationConfigMap(ctx context.Context, c client.Client, namespace, name, main, variables string) (*corev1.ConfigMap, error) {
	return createOrUpdateConfigMap(ctx, c, namespace, name, map[string]string{
		MainKey:      main,
		VariablesKey: variables,
	})
}

func CreateOrUpdateStateConfigMap(ctx context.Context, c client.Client, namespace, name, state string) (*corev1.ConfigMap, error) {
	return createOrUpdateConfigMap(ctx, c, namespace, name, map[string]string{
		StateKey: state,
	})
}

func CreateOrUpdateTFVarsSecret(ctx context.Context, c client.Client, namespace, name string, tfvars []byte) (*corev1.Secret, error) {
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}
	err := kutil.CreateOrUpdate(ctx, c, secret, func() error {
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		secret.Data[TFVarsKey] = tfvars
		return nil
	})

	if err != nil {
		return nil, err
	}
	return secret, nil
}

// prepare checks whether all required ConfigMaps and Secrets exist. It returns the number of
// existing ConfigMaps/Secrets, or the error in case something unexpected happens.
func (t *Terraformer) prepare(ctx context.Context) (int, error) {
	numberOfExistingResources, err := t.verifyConfigExists(ctx)
	if err != nil {
		return -1, err
	}

	if t.variablesEnvironment == nil {
		return -1, errors.New("no Terraform variables environment provided")
	}

	// Clean up possible existing job/pod artifacts from previous runs
	if err := t.EnsureCleanedUp(); err != nil {
		return -1, err
	}

	return numberOfExistingResources, nil
}

func (t *Terraformer) verifyConfigExists(ctx context.Context) (int, error) {
	numberOfExistingResources := 0

	if err := t.client.Get(ctx, kutil.Key(t.namespace, t.stateName), &corev1.ConfigMap{}); err == nil {
		numberOfExistingResources++
	} else if err != nil && !apierrors.IsNotFound(err) {
		return -1, err
	}

	if err := t.client.Get(ctx, kutil.Key(t.namespace, t.variablesName), &corev1.Secret{}); err == nil {
		numberOfExistingResources++
	} else if err != nil && !apierrors.IsNotFound(err) {
		return -1, err
	}

	if err := t.client.Get(ctx, kutil.Key(t.namespace, t.configName), &corev1.ConfigMap{}); err == nil {
		numberOfExistingResources++
	} else if err != nil && !apierrors.IsNotFound(err) {
		return -1, err
	}

	return numberOfExistingResources, nil
}

// ConfigExists returns true if all three Terraform configuration secrets/configmaps exist, and false otherwise.
func (t *Terraformer) ConfigExists() (bool, error) {
	numberOfExistingResources, err := t.verifyConfigExists(context.TODO())
	return numberOfExistingResources == numberOfConfigResources, err
}

// cleanupConfiguration deletes the two ConfigMaps which store the Terraform configuration and state. It also deletes
// the Secret which stores the Terraform variables.
func (t *Terraformer) cleanupConfiguration(ctx context.Context) error {
	t.logger.Debugf("Deleting Terraform variables Secret '%s'", t.variablesName)
	if err := t.client.Delete(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: t.namespace, Name: t.variablesName}}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	t.logger.Debugf("Deleting Terraform configuration ConfigMap '%s'", t.configName)
	if err := t.client.Delete(ctx, &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: t.namespace, Name: t.configName}}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	t.logger.Debugf("Deleting Terraform state ConfigMap '%s'", t.stateName)
	if err := t.client.Delete(ctx, &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: t.namespace, Name: t.stateName}}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

// EnsureCleanedUp deletes the job, pods, and waits until everything has been cleaned up.
func (t *Terraformer) EnsureCleanedUp() error {
	ctx := context.TODO()
	jobPodList, err := t.listJobPods(ctx)
	if err != nil {
		return err
	}
	if err := t.cleanupJob(ctx, jobPodList); err != nil {
		return err
	}
	return t.WaitForCleanEnvironment(ctx)
}
