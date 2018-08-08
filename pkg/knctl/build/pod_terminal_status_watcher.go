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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type PodTerminalStatusWatcher struct {
	Pod        corev1.Pod
	PodsClient typedcorev1.PodInterface
}

func (l PodTerminalStatusWatcher) Wait(cancelCh chan struct{}) (corev1.PodPhase, error) {
	for {
		// TODO infinite retry

		pod, err := l.PodsClient.Get(l.Pod.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			return pod.Status.Phase, nil
		}

		select {
		case <-cancelCh:
			return pod.Status.Phase, nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
