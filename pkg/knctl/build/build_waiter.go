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
	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type BuildWaiter struct {
	build            *v1alpha1.Build
	buildClient      buildclientset.Interface
	podsGetterClient typedcorev1.PodsGetter
}

func NewBuildWaiter(
	build *v1alpha1.Build,
	buildClient buildclientset.Interface,
	podsGetterClient typedcorev1.PodsGetter,
) BuildWaiter {
	return BuildWaiter{build, buildClient, podsGetterClient}
}

func (w BuildWaiter) WaitForBuilderAssignment(cancelCh chan struct{}) (*v1alpha1.Build, error) {
	for {
		// TODO infinite retry

		build, err := w.buildClient.BuildV1alpha1().Builds(w.build.Namespace).Get(w.build.Name, metav1.GetOptions{})
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

		build, err := w.buildClient.BuildV1alpha1().Builds(w.build.Namespace).Get(w.build.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Getting build while waiting for completion: %s", err)
		}

		cond := build.Status.GetCondition(v1alpha1.BuildSucceeded)
		if cond != nil {
			switch cond.Status {
			case corev1.ConditionTrue, corev1.ConditionFalse:
				return build, nil
			default:
				// continue waiting
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

func (w BuildWaiter) WaitForClusterBuilderPodAssignment(cancelCh chan struct{}) (*v1alpha1.Build, *corev1.Pod, error) {
	var build *v1alpha1.Build

	for {
		// TODO infinite retry

		var err error

		build, err = w.buildClient.BuildV1alpha1().Builds(w.build.Namespace).Get(w.build.Name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("Getting build while waiting for cluster build to assign pod: %s", err)
		}

		if build.Status.Cluster != nil {
			if len(build.Status.Cluster.Namespace) > 0 && len(build.Status.Cluster.PodName) > 0 {
				break
			}
		}

		select {
		case <-cancelCh:
			return build, nil, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}

	// Check if pod was initialized and is ready to be interacted via the API
	for {
		// TODO infinite retry
		podsClient := w.podsGetterClient.Pods(build.Status.Cluster.Namespace)

		pod, err := podsClient.Get(build.Status.Cluster.PodName, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return build, nil, fmt.Errorf("Getting assigned building pod: %s", err)
			}
		} else {
			return build, pod, nil
		}

		select {
		case <-cancelCh:
			return build, pod, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
