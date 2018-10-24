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

package service

import (
	"fmt"

	ctling "github.com/cppforlife/knctl/pkg/knctl/ingress"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type ServiceAddress struct {
	service    *v1alpha1.Service
	coreClient kubernetes.Interface
}

func (o ServiceAddress) Domain() (string, error) {
	if len(o.service.Status.Domain) == 0 {
		return "", fmt.Errorf("Expected service '%s' to have non-empty domain", o.service.Name)
	}

	return o.service.Status.Domain, nil
}

func (o ServiceAddress) URL(port int32, useDomain bool) (string, error) {
	domain, err := o.Domain()
	if err != nil {
		return "", err
	}

	ingressAddress, ingressPort, err := ctling.NewIngressServices(o.coreClient).PreferredAddress(port)
	if err != nil {
		return "", err
	}

	host := ingressAddress

	if useDomain {
		host = domain
	}

	return o.schema(port) + "://" + host + ":" + ingressPort, nil
}

func (o ServiceAddress) schema(port int32) string {
	if port == 443 {
		return "https"
	}
	return "http"
}
