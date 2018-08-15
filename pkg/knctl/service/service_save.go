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

package service

import (
	"fmt"

	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) CreateOrUpdate() (*v1alpha1.Service, error) {
	createdService, createErr := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Create(&s.service)
	if createErr != nil {
		if errors.IsAlreadyExists(createErr) {
			origService, err := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Get(s.service.Name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("Creating service: %s", createErr)
			}

			// TODO currently replaces everything
			s.service.ResourceVersion = origService.ResourceVersion

			updatedService, updateErr := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Update(&s.service)
			if updateErr != nil {
				return nil, fmt.Errorf("Updating service: %s", updateErr)
			}

			s.service = *updatedService

			return updatedService, nil
		}

		return nil, fmt.Errorf("Creating service: %s", createErr)
	}

	s.service = *createdService

	return createdService, nil
}
