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
	"sync"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/cppforlife/knctl/pkg/knctl/logs"
	"github.com/knative/serving/pkg/apis/serving"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type podWatcher interface {
	Watch(podsToWatchCh chan corev1.Pod, cancelCh chan struct{}) error
}

type LogsView struct {
	tailOpts      logs.PodLogOpts
	podWatcher    podWatcher
	coreClient    kubernetes.Interface
	ui            ui.UI
	cancelSignals cmdcore.CancelSignals
}

func (v LogsView) Show() error {
	podsToWatchCh, cancelPodTailCh, err := v.setUpPodWatching()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for pod := range podsToWatchCh {
		pod := pod
		wg.Add(1)

		go func() {
			podsClient := v.coreClient.CoreV1().Pods(pod.Namespace)
			tag := fmt.Sprintf("%s > %s", pod.Labels[serving.RevisionLabelKey], pod.Name)

			err := logs.NewPodContainerLog(pod, "user-container", podsClient, tag, v.tailOpts).Tail(v.ui, cancelPodTailCh)
			if err != nil {
				v.ui.BeginLinef("Pod logs tailing error: %s\n", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}

func (v LogsView) setUpPodWatching() (chan corev1.Pod, chan struct{}, error) {
	podsToWatchCh := make(chan corev1.Pod)
	cancelPodTailCh := make(chan struct{})
	cancelCh := make(chan struct{})

	if v.tailOpts.Follow {
		v.cancelSignals.Watch(func() {
			close(cancelCh)
			close(cancelPodTailCh)
		})
	} else {
		close(cancelCh)
	}

	go func() {
		v.podWatcher.Watch(podsToWatchCh, cancelCh)
		close(podsToWatchCh)
	}()

	return podsToWatchCh, cancelPodTailCh, nil
}
