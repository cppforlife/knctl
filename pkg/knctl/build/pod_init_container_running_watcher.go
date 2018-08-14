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

type PodInitContainerRunningWatcher struct {
	Pod           corev1.Pod
	PodsClient    typedcorev1.PodInterface
	InitContainer string
}

func (l PodInitContainerRunningWatcher) Wait(cancelCh chan struct{}) error {
	for {
		// TODO infinite retry

		pod, err := l.PodsClient.Get(l.Pod.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		for _, status := range pod.Status.InitContainerStatuses {
			if status.Name == l.InitContainer {
				// TODO what if pod is no longer progressing?
				if status.State.Running != nil {
					return nil
				}
			}
		}

		select {
		case <-cancelCh:
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
