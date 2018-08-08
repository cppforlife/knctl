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

package logs

import (
	"sync"

	"github.com/cppforlife/go-cli-ui/ui"
	corev1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type PodLogOpts struct {
	Follow bool
	Lines  *int64
}

type PodLog struct {
	pod        corev1.Pod
	podsClient typedcorev1.PodInterface

	tag  string
	opts PodLogOpts
}

func NewPodLog(
	pod corev1.Pod,
	podsClient typedcorev1.PodInterface,
	tag string,
	opts PodLogOpts,
) PodLog {
	return PodLog{pod, podsClient, tag, opts}
}

// TailAll will tail all logs from all containers in a single Pod
func (l PodLog) TailAll(ui ui.UI, cancelCh chan struct{}) error {
	var wg sync.WaitGroup

	var conts []corev1.Container
	conts = append(conts, l.pod.Spec.InitContainers...)
	conts = append(conts, l.pod.Spec.Containers...)

	for _, cont := range conts {
		cont := cont
		wg.Add(1)

		go func() {
			NewPodContainerLog(l.pod, cont.Name, l.podsClient, l.tag, l.opts).Tail(ui, cancelCh) // TODO err?
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
