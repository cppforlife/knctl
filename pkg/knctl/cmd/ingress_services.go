/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Open 2.0 (the "License");
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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type IngressServices struct {
	coreClient kubernetes.Interface
}

func (s IngressServices) List() (*corev1.ServiceList, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			"knative": "ingressgateway",
		}).String(),
	}

	istioNsName := NewIstio(s.coreClient).SystemNamespaceName()

	services, err := s.coreClient.CoreV1().Services(istioNsName).List(listOpts)
	if err != nil {
		return nil, fmt.Errorf("Listing services in istio namespace: %s", err)
	}

	return services, nil
}
