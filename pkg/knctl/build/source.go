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
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type SourceFactory struct {
	waiter     BuildWaiter
	coreClient kubernetes.Interface
	restConfig *rest.Config
}

type Source interface {
	Upload(ui ui.UI, cancelCh chan struct{}) error
}

func NewSourceFactory(
	waiter BuildWaiter,
	coreClient kubernetes.Interface,
	restConfig *rest.Config,
) SourceFactory {
	return SourceFactory{waiter, coreClient, restConfig}
}

func (s SourceFactory) New(t v1alpha1.BuildProvider, dirPath string) Source {
	switch t {
	case v1alpha1.ClusterBuildProvider:
		return NewClusterBuilderSource(dirPath, s.waiter, s.coreClient, s.restConfig)

	default:
		return NoopSource{t}
	}
}
