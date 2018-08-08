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
	"time"

	"github.com/knative/build/pkg/apis/build/v1alpha1"
	typedv1alpha1 "github.com/knative/build/pkg/client/clientset/versioned/typed/build/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildWaiter struct {
	build        *v1alpha1.Build
	buildsClient typedv1alpha1.BuildInterface
}

func NewBuildWaiter(build *v1alpha1.Build, buildsClient typedv1alpha1.BuildInterface) BuildWaiter {
	return BuildWaiter{build, buildsClient}
}

func (w BuildWaiter) WaitForBuilderAssignment(cancelCh chan struct{}) (*v1alpha1.Build, error) {
	for {
		// TODO infinite retry

		build, err := w.buildsClient.Get(w.build.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Getting build while waiting for builder assignment: %s", err)
		}

		if len(build.Status.Builder) > 0 {
			return build, nil
		}

		select {
		case <-cancelCh:
			return build, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (w BuildWaiter) WaitForCompletion(cancelCh chan struct{}) (*v1alpha1.Build, error) {
	for {
		// TODO infinite retry

		build, err := w.buildsClient.Get(w.build.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Getting build while waiting for completion: %s", err)
		}

		cond := build.Status.GetCondition(v1alpha1.BuildSucceeded)
		if cond != nil {
			return build, nil
		}

		select {
		case <-cancelCh:
			return build, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (w BuildWaiter) WaitForClusterBuilderPodAssignment(cancelCh chan struct{}) (*v1alpha1.Build, error) {
	for {
		// TODO infinite retry

		build, err := w.buildsClient.Get(w.build.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Getting build while waiting for cluster build to assign pod: %s", err)
		}

		if build.Status.Cluster != nil {
			if len(build.Status.Cluster.Namespace) > 0 && len(build.Status.Cluster.PodName) > 0 {
				return build, nil
			}
		}

		select {
		case <-cancelCh:
			return build, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
