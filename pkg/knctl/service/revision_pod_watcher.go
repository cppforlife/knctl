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

package service

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type RevisionPodWatcher struct {
	revision *v1alpha1.Revision

	servingClient servingclientset.Interface
	coreClient    kubernetes.Interface

	ui ui.UI // TODO remove
}

func NewRevisionPodWatcher(
	revision *v1alpha1.Revision,
	servingClient servingclientset.Interface,
	coreClient kubernetes.Interface,
	ui ui.UI,
) RevisionPodWatcher {
	return RevisionPodWatcher{revision, servingClient, coreClient, ui}
}

func (w RevisionPodWatcher) Watch(podsToWatchCh chan corev1.Pod, cancelCh chan struct{}) error {
	nonUniquePodsToWatchCh := make(chan corev1.Pod)

	go func() {
		podWatcher := NewPodWatcher(
			w.coreClient.CoreV1().Pods(w.revision.Namespace),
			metav1.ListOptions{
				LabelSelector: labels.Set(map[string]string{
					serving.RevisionUID: string(w.revision.UID),
				}).String(),
			},
		)

		err := podWatcher.Watch(nonUniquePodsToWatchCh, cancelCh)
		if err != nil {
			w.ui.BeginLinef("Pod watching error: %s\n", err)
		}

		close(nonUniquePodsToWatchCh)
	}()

	// Send unique pods to the watcher client
	watchedPods := map[string]struct{}{}

	for pod := range nonUniquePodsToWatchCh {
		podUID := string(pod.UID)
		if _, found := watchedPods[podUID]; found {
			continue
		}

		watchedPods[podUID] = struct{}{}
		podsToWatchCh <- pod
	}

	return nil
}
