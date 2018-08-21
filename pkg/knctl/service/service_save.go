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
	"time"

	"github.com/cppforlife/knctl/pkg/knctl/util"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) CreateOrUpdate() (*v1alpha1.Service, error) {
	createdService, createErr := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Create(&s.service)
	if createErr != nil {
		if errors.IsAlreadyExists(createErr) {
			return s.update()
		}

		return nil, fmt.Errorf("Creating service: %s", createErr)
	}

	s.service = *createdService

	return createdService, nil
}

func (s *Service) update() (*v1alpha1.Service, error) {
	var updatedService *v1alpha1.Service

	updateErr := util.Retry(time.Second, 10*time.Second, func() (bool, error) {
		origService, err := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Get(s.service.Name, metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("Creating service: %s", err)
		}

		origService.Spec = s.service.Spec

		service, err := s.servingClient.ServingV1alpha1().Services(s.service.Namespace).Update(origService)
		if err != nil {
			return false, fmt.Errorf("Updating service: %s", err)
		}

		updatedService = service

		return true, nil
	})
	if updateErr != nil {
		return nil, updateErr
	}

	s.service = *updatedService

	return updatedService, nil
}
