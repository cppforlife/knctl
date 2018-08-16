/*
Copyright 2018 The Knative Authors

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

package cmd

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	istioNsInjectionLabel   = "istio-injection"
	istioNsInjectionEnabled = "enabled"
)

type Istio struct {
	coreClient kubernetes.Interface
}

func NewIstio(coreClient kubernetes.Interface) Istio {
	return Istio{coreClient}
}

func (i Istio) SystemNamespaceName() string { return "istio-system" }

func (i Istio) EnableNamespace(namespaceName string) error {
	ns, err := i.coreClient.CoreV1().Namespaces().Get(namespaceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if ns.Labels == nil {
		ns.Labels = map[string]string{}
	}

	if _, found := ns.Labels[istioNsInjectionLabel]; found {
		return nil
	}

	ns.Labels[istioNsInjectionLabel] = istioNsInjectionEnabled

	_, err = i.coreClient.CoreV1().Namespaces().Update(ns)
	if err != nil {
		return fmt.Errorf("Updating namespace: %s", err)
	}

	return nil
}
