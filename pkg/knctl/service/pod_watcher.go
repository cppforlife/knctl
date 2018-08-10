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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type PodWatcher struct {
	podsClient typedcorev1.PodInterface
	listOpts   metav1.ListOptions
}

func NewPodWatcher(
	podsClient typedcorev1.PodInterface,
	listOpts metav1.ListOptions,
) PodWatcher {
	return PodWatcher{podsClient, listOpts}
}

func (w PodWatcher) Watch(podsToWatchCh chan corev1.Pod, cancelCh chan struct{}) error {
	podsList, err := w.podsClient.List(w.listOpts)
	if err != nil {
		return err
	}

	for _, pod := range podsList.Items {
		podsToWatchCh <- pod
	}

	// Return before potentially getting any events
	select {
	case <-cancelCh:
		return nil
	default:
	}

	for {
		retry, err := w.watch(podsToWatchCh, cancelCh)
		if err != nil {
			return err
		}
		if !retry {
			return nil
		}
	}
}

func (w PodWatcher) watch(podsToWatchCh chan corev1.Pod, cancelCh chan struct{}) (bool, error) {
	watcher, err := w.podsClient.Watch(w.listOpts)
	if err != nil {
		return false, fmt.Errorf("Creating Pod watcher: %s", err)
	}

	defer watcher.Stop()

	for {
		select {
		case e, ok := <-watcher.ResultChan():
			if !ok || e.Object == nil {
				// Watcher may expire, hence try to retry
				return true, nil
			}

			pod, ok := e.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			switch e.Type {
			case watch.Added:
				podsToWatchCh <- *pod
			}

		case <-cancelCh:
			return false, nil
		}
	}
}
