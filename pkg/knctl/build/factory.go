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

package build

import (
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Factory struct {
	buildClient buildclientset.Interface
	coreClient  kubernetes.Interface
	restConfig  *rest.Config
}

func NewFactory(
	buildClient buildclientset.Interface,
	coreClient kubernetes.Interface,
	restConfig *rest.Config,
) Factory {
	return Factory{buildClient, coreClient, restConfig}
}

func (f Factory) New(build *v1alpha1.Build) Build {
	waiter := NewBuildWaiter(build, f.buildClient, f.coreClient.CoreV1())
	logs := NewLogs(waiter, f.coreClient.CoreV1())
	sourceFactory := NewSourceFactory(waiter, f.coreClient, f.restConfig)
	return NewBuild(waiter, logs, sourceFactory)
}
