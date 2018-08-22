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
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Logs struct {
	waiter           BuildWaiter
	podsGetterClient typedcorev1.PodsGetter
}

func NewLogs(waiter BuildWaiter, podsGetterClient typedcorev1.PodsGetter) Logs {
	return Logs{waiter, podsGetterClient}
}

func (l Logs) Tail(ui ui.UI, cancelCh chan struct{}) error {
	ui.PrintLinef("Watching build logs...")

	build, err := l.waiter.WaitForBuilderAssignment(cancelCh)
	if err != nil {
		return fmt.Errorf("Waiting for build to be assigned a builder: %s", err)
	}

	switch build.Status.Builder {
	case v1alpha1.ClusterBuildProvider:
		return NewClusterBuilderLogs(l.waiter, l.podsGetterClient).Tail(ui, cancelCh)

	default:
		ui.PrintLinef("Cannot follow logs unknown builder '%s'...\n", build.Status.Builder)
		return nil
	}
}
