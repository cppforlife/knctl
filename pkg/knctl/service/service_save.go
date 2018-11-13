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
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) CreateOrUpdate() (*v1alpha1.Service, error) {
	service, err := s.serviceSpec.Service()
	if err != nil {
		return nil, err
	}

	conf, err := s.serviceSpec.Configuration()
	if err != nil {
		return nil, err
	}

	createdService, err := s.createOrUpdateService(service)
	if err != nil {
		return nil, err
	}

	if s.serviceSpec.NeedsConfigurationUpdate() {
		_, err = s.createOrUpdateConfiguration(createdService, conf)
		if err != nil {
			return nil, err
		}
	}

	return createdService, nil
}

func (s *Service) createOrUpdateService(service v1alpha1.Service) (*v1alpha1.Service, error) {
	createdService, createErr := s.servingClient.ServingV1alpha1().Services(s.serviceSpec.Namespace()).Create(&service)
	if createErr != nil {
		if errors.IsAlreadyExists(createErr) {
			return s.updateService(service)
		}

		return nil, fmt.Errorf("Creating service: %s", createErr)
	}

	return createdService, nil
}

func (s *Service) updateService(service v1alpha1.Service) (*v1alpha1.Service, error) {
	var updatedService *v1alpha1.Service

	updateErr := util.Retry(time.Second, 10*time.Second, func() (bool, error) {
		origService, err := s.servingClient.ServingV1alpha1().Services(s.serviceSpec.Namespace()).Get(s.serviceSpec.Name(), metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("Creating service: %s", err)
		}

		origService.Spec = service.Spec

		service, err := s.servingClient.ServingV1alpha1().Services(s.serviceSpec.Namespace()).Update(origService)
		if err != nil {
			return false, fmt.Errorf("Updating service: %s", err)
		}

		updatedService = service

		return true, nil
	})
	if updateErr != nil {
		return nil, updateErr
	}

	return updatedService, nil
}

func (s *Service) createOrUpdateConfiguration(service *v1alpha1.Service, conf v1alpha1.Configuration) (*v1alpha1.Configuration, error) {
	conf.ObjectMeta = metav1.ObjectMeta{
		Name:      service.Name,
		Namespace: service.Namespace,
		Labels:    s.serviceLabels(service),
		// Configure configuration to be deleted when service is deleted
		OwnerReferences: []metav1.OwnerReference{
			*metav1.NewControllerRef(service.GetObjectMeta(), service.GetGroupVersionKind()),
		},
	}

	createdConfiguration, createErr := s.servingClient.ServingV1alpha1().Configurations(s.serviceSpec.Namespace()).Create(&conf)
	if createErr != nil {
		if errors.IsAlreadyExists(createErr) {
			return s.updateConfiguration(conf)
		}

		return nil, fmt.Errorf("Creating conf: %s", createErr)
	}

	return createdConfiguration, nil
}

func (s *Service) updateConfiguration(conf v1alpha1.Configuration) (*v1alpha1.Configuration, error) {
	var updatedConfiguration *v1alpha1.Configuration

	updateErr := util.Retry(time.Second, 10*time.Second, func() (bool, error) {
		origConfiguration, err := s.servingClient.ServingV1alpha1().Configurations(s.serviceSpec.Namespace()).Get(s.serviceSpec.Name(), metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("Creating conf: %s", err)
		}

		conf.Spec.Generation = origConfiguration.Spec.Generation
		origConfiguration.Spec = conf.Spec

		conf, err := s.servingClient.ServingV1alpha1().Configurations(s.serviceSpec.Namespace()).Update(origConfiguration)
		if err != nil {
			return false, fmt.Errorf("Updating conf: %s", err)
		}

		updatedConfiguration = conf

		return true, nil
	})
	if updateErr != nil {
		return nil, updateErr
	}

	return updatedConfiguration, nil
}

// Taken from https://github.com/knative/serving/blob/master/pkg/reconciler/v1alpha1/service/resources/labels.go
func (*Service) serviceLabels(s *v1alpha1.Service) map[string]string {
	labels := make(map[string]string, len(s.ObjectMeta.Labels)+1)
	labels[serving.ServiceLabelKey] = s.Name

	// Pass through the labels on the Service to child resources.
	for k, v := range s.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}
