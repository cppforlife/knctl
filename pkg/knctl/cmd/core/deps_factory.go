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

package core

import (
	"fmt"

	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type DepsFactory interface {
	ServingClient() (servingclientset.Interface, error)
	BuildClient() (buildclientset.Interface, error)
	CoreClient() (kubernetes.Interface, error)
}

func NewDepsFactory() DepsFactory { // Concise for testing
	return NewDepsFactoryImpl(NewConfigFactoryImpl())
}

type DepsFactoryImpl struct {
	configFactory ConfigFactory
}

var _ DepsFactory = &DepsFactoryImpl{}

func NewDepsFactoryImpl(configFactory ConfigFactory) *DepsFactoryImpl {
	return &DepsFactoryImpl{configFactory}
}

func (f *DepsFactoryImpl) ServingClient() (servingclientset.Interface, error) {
	config, err := f.configFactory.RESTConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := servingclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Building Serving clientset: %s", err)
	}

	return clientset, nil
}

func (f *DepsFactoryImpl) BuildClient() (buildclientset.Interface, error) {
	config, err := f.configFactory.RESTConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := buildclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Building Build clientset: %s", err)
	}

	return clientset, nil
}

func (f *DepsFactoryImpl) CoreClient() (kubernetes.Interface, error) {
	config, err := f.configFactory.RESTConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Building Core clientset: %s", err)
	}

	return clientset, nil
}
