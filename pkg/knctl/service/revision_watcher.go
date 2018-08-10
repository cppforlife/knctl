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

	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	typedv1alpha1 "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type RevisionWatcher struct {
	revisionsClient typedv1alpha1.RevisionInterface
	listOpts        metav1.ListOptions
}

func NewRevisionWatcher(
	revisionsClient typedv1alpha1.RevisionInterface,
	listOpts metav1.ListOptions,
) RevisionWatcher {
	return RevisionWatcher{revisionsClient, listOpts}
}

func (w RevisionWatcher) Watch(revisionsToWatchCh chan v1alpha1.Revision, cancelCh chan struct{}) error {
	revisionsList, err := w.revisionsClient.List(w.listOpts)
	if err != nil {
		return err
	}

	for _, revision := range revisionsList.Items {
		revisionsToWatchCh <- revision
	}

	// Return before potentially getting any events
	select {
	case <-cancelCh:
		return nil
	default:
	}

	for {
		retry, err := w.watch(revisionsToWatchCh, cancelCh)
		if err != nil {
			return err
		}
		if !retry {
			return nil
		}
	}
}

func (w RevisionWatcher) watch(revisionsToWatchCh chan v1alpha1.Revision, cancelCh chan struct{}) (bool, error) {
	watcher, err := w.revisionsClient.Watch(w.listOpts)
	if err != nil {
		return false, fmt.Errorf("Creating Revision watcher: %s", err)
	}

	defer watcher.Stop()

	for {
		select {
		case e, ok := <-watcher.ResultChan():
			if !ok || e.Object == nil {
				// Watcher may expire, hence try to retry
				return true, nil
			}

			revision, ok := e.Object.(*v1alpha1.Revision)
			if !ok {
				continue
			}

			switch e.Type {
			case watch.Added:
				revisionsToWatchCh <- *revision
			}

		case <-cancelCh:
			return false, nil
		}
	}
}
