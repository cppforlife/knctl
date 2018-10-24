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

package domain

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	domainsNs            = "knative-serving"
	domainsConfigMapName = "config-domain"
	domainsDefaultValue  = ""
)

type Domains struct {
	coreClient kubernetes.Interface
}

type Domain struct {
	Name    string
	Default bool
}

func NewDomains(coreClient kubernetes.Interface) Domains {
	return Domains{coreClient}
}

func (d Domains) List() ([]Domain, error) {
	config, err := d.coreClient.CoreV1().ConfigMaps(domainsNs).Get(domainsConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var domains []Domain

	for k, v := range config.Data {
		// Empty value indicates default domain
		domains = append(domains, Domain{
			Name:    k,
			Default: v == domainsDefaultValue,
		})
	}

	return domains, nil
}

func (d Domains) Create(domain Domain) error {
	config, err := d.coreClient.CoreV1().ConfigMaps(domainsNs).Get(domainsConfigMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if domain.Default {
		for k, v := range config.Data {
			if v == domainsDefaultValue {
				delete(config.Data, k)
				break
			}
		}

		config.Data[domain.Name] = domainsDefaultValue
	} else {
		return fmt.Errorf("Setting non-default domain is not supported")
	}

	_, err = d.coreClient.CoreV1().ConfigMaps(domainsNs).Update(config)
	if err != nil {
		return fmt.Errorf("Updating domains: %s", err)
	}

	return nil
}
