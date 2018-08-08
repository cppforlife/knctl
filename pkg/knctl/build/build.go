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

	"github.com/cppforlife/go-cli-ui/ui" // TODO replace
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	typedv1alpha1 "github.com/knative/build/pkg/client/clientset/versioned/typed/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Build struct {
	build            *v1alpha1.Build
	buildsClient     typedv1alpha1.BuildInterface
	podsGetterClient typedcorev1.PodsGetter
}

func NewBuild(
	build *v1alpha1.Build,
	buildsClient typedv1alpha1.BuildInterface,
	podsGetterClient typedcorev1.PodsGetter,
) Build {
	return Build{build, buildsClient, podsGetterClient}
}

func (b Build) TailLogs(ui ui.UI, cancelCh chan struct{}) error {
	return NewLogs(b.build, b.buildsClient, b.podsGetterClient).Tail(ui, cancelCh)
}

func (b Build) Error(cancelCh chan struct{}) error {
	build, err := NewBuildWaiter(b.build, b.buildsClient).WaitForCompletion(cancelCh)
	if err != nil {
		return err
	}

	cond := build.Status.GetCondition(v1alpha1.BuildSucceeded)
	if cond == nil {
		return fmt.Errorf("Expected build to complete")
	}

	switch cond.Status {
	case corev1.ConditionTrue:
		return nil
	case corev1.ConditionFalse:
		return fmt.Errorf("Build failed")
	default:
		return fmt.Errorf("Build may or may not have completed (state 'Unknown')")
	}
}
