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

package route

import (
	"fmt"

	ctling "github.com/cppforlife/knctl/pkg/knctl/ingress"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// TODO consolidate with ServiceAddress
type RouteAddress struct {
	route      *v1alpha1.Route
	coreClient kubernetes.Interface
}

func (o RouteAddress) Domain() (string, error) {
	if len(o.route.Status.Domain) == 0 {
		return "", fmt.Errorf("Expected route '%s' to have non-empty domain", o.route.Name)
	}

	return o.route.Status.Domain, nil
}

func (o RouteAddress) URL(port int32, useDomain bool) (string, error) {
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

func (o RouteAddress) schema(port int32) string {
	if port == 443 {
		return "https"
	}
	return "http"
}
