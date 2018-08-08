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
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	typedv1alpha1 "github.com/knative/build/pkg/client/clientset/versioned/typed/build/v1alpha1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Logs struct {
	build            *v1alpha1.Build
	buildsClient     typedv1alpha1.BuildInterface
	podsGetterClient typedcorev1.PodsGetter
}

func NewLogs(
	build *v1alpha1.Build,
	buildsClient typedv1alpha1.BuildInterface,
	podsGetterClient typedcorev1.PodsGetter,
) Logs {
	return Logs{build, buildsClient, podsGetterClient}
}

func (l Logs) Tail(ui ui.UI, cancelCh chan struct{}) error {
	// TODO no new build is kicked off on subsequent deploys?
	ui.PrintLinef("Watching build logs...")

	build, err := NewBuildWaiter(l.build, l.buildsClient).WaitForBuilderAssignment(cancelCh)
	if err != nil {
		return fmt.Errorf("Waiting for build to be assigned a builder: %s", err)
	}

	switch build.Status.Builder {
	case v1alpha1.ClusterBuildProvider:
		return NewClusterBuilderLogs(build, l.buildsClient, l.podsGetterClient).Tail(ui, cancelCh)
	case v1alpha1.GoogleBuildProvider:
		return NewGoogleBuilderLogs().Tail(ui, cancelCh)
	default:
		ui.PrintLinef("Cannot follow logs unknown builder '%s'...\n", build.Status.Builder)
		return nil
	}
}
