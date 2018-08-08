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

	"github.com/knative/build/pkg/apis/build/v1alpha1"
	typedv1alpha1 "github.com/knative/build/pkg/client/clientset/versioned/typed/build/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type BuildWatcher struct {
	buildsClient typedv1alpha1.BuildInterface
	listOpts     metav1.ListOptions
}

func NewBuildWatcher(
	buildsClient typedv1alpha1.BuildInterface,
	listOpts metav1.ListOptions,
) BuildWatcher {
	return BuildWatcher{buildsClient, listOpts}
}

func (w BuildWatcher) Watch(buildsToWatch chan v1alpha1.Build, cancelCh chan struct{}) error {
	watcher, err := w.buildsClient.Watch(w.listOpts)
	if err != nil {
		return fmt.Errorf("Creating Build watcher: %s", err)
	}

	defer watcher.Stop()

	buildsList, err := w.buildsClient.List(w.listOpts)
	if err != nil {
		return err
	}

	for _, build := range buildsList.Items {
		buildsToWatch <- build
	}

	for {
		select {
		case e := <-watcher.ResultChan():
			if e.Object == nil {
				return nil // TODO return?
			}

			build, ok := e.Object.(*v1alpha1.Build)
			if !ok {
				continue
			}

			switch e.Type {
			case watch.Added:
				buildsToWatch <- *build
			}

		case <-cancelCh:
			return nil
		}
	}
}
