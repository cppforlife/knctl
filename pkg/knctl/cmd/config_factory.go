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

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigFactory interface {
	ConfigurePathResolver(func() (string, error))
	RESTConfig() (*rest.Config, error)
	DefaultNamespace() (string, error)
}

type ConfigFactoryImpl struct {
	pathResolverFunc func() (string, error)
}

var _ ConfigFactory = &ConfigFactoryImpl{}

func NewConfigFactoryImpl() *ConfigFactoryImpl {
	return &ConfigFactoryImpl{}
}

func (f *ConfigFactoryImpl) ConfigurePathResolver(resolverFunc func() (string, error)) {
	f.pathResolverFunc = resolverFunc
}

func (f *ConfigFactoryImpl) RESTConfig() (*rest.Config, error) {
	config, err := f.clientConfig()
	if err != nil {
		return nil, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Building Kubernetes config: %s", err)
	}

	return restConfig, nil
}

func (f *ConfigFactoryImpl) DefaultNamespace() (string, error) {
	config, err := f.clientConfig()
	if err != nil {
		return "", err
	}

	name, _, err := config.Namespace()
	return name, err
}

func (f *ConfigFactoryImpl) clientConfig() (clientcmd.ClientConfig, error) {
	path, err := f.pathResolverFunc()
	if err != nil {
		return nil, fmt.Errorf("Resolving config path: %s", err)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{},
	), nil
}
