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
	ConfigurePath(string)
	RESTConfig() (*rest.Config, error)
	DefaultNamespace() (string, error)
}

type ConfigFactoryImpl struct {
	configPath string
}

var _ ConfigFactory = &ConfigFactoryImpl{}

func NewConfigFactoryImpl() *ConfigFactoryImpl {
	return &ConfigFactoryImpl{}
}

func (f *ConfigFactoryImpl) ConfigurePath(path string) {
	f.configPath = path
}

func (f *ConfigFactoryImpl) RESTConfig() (*rest.Config, error) {
	config, err := f.clientConfig().ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Building Kubernetes config: %s", err)
	}

	return config, nil
}

func (f *ConfigFactoryImpl) DefaultNamespace() (string, error) {
	name, _, err := f.clientConfig().Namespace()
	return name, err
}

func (f *ConfigFactoryImpl) clientConfig() clientcmd.ClientConfig {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: f.configPath},
		&clientcmd.ConfigOverrides{},
	)
}
