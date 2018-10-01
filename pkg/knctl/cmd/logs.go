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

package cmd

import (
	"fmt"
	"sync"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/cppforlife/knctl/pkg/knctl/logs"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LogsOptions struct {
	ui            ui.UI
	depsFactory   DepsFactory
	cancelSignals CancelSignals

	ServiceFlags ServiceFlags

	Follow bool
	Lines  int64
}

func NewLogsOptions(ui ui.UI, depsFactory DepsFactory, cancelSignals CancelSignals) *LogsOptions {
	return &LogsOptions{ui: ui, depsFactory: depsFactory, cancelSignals: cancelSignals}
}

func NewLogsCmd(o *LogsOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Print service logs",
		Long:  "Print service logs of all active pods for a service",
		Example: `
  # Fetch last 10 log lines for service 'svc1' in namespace 'ns1' 
  knctl logs -s svc1 -n ns1

  # Follow logs for service 'svc1' in namespace 'ns1' 
  knctl logs -f -s svc1 -n ns1`,
		Annotations: map[string]string{
			basicGroup.Key: basicGroup.Value,
		},
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)

	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", false, "As new revisions are added, new pod logs will be printed")
	cmd.Flags().Int64VarP(&o.Lines, "lines", "l", 10, "Number of lines")

	return cmd
}

func (o *LogsOptions) Run() error {
	if !o.Follow && o.Lines <= 0 {
		return fmt.Errorf("Expected --lines to be greater than zero since --follow is not specified")
	}

	tailOpts := logs.PodLogOpts{Follow: o.Follow}

	if !o.Follow {
		tailOpts.Lines = &o.Lines
	}

	podsToWatchCh, cancelPodTailCh, err := o.setUpPodWatching()
	if err != nil {
		return err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for pod := range podsToWatchCh {
		pod := pod
		wg.Add(1)

		go func() {
			podsClient := coreClient.CoreV1().Pods(pod.Namespace)
			tag := fmt.Sprintf("%s > %s", pod.Labels[serving.RevisionLabelKey], pod.Name)

			err := logs.NewPodContainerLog(pod, "user-container", podsClient, tag, tailOpts).Tail(o.ui, cancelPodTailCh)
			if err != nil {
				o.ui.BeginLinef("Pod logs tailing error: %s\n", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}

func (o *LogsOptions) setUpPodWatching() (chan corev1.Pod, chan struct{}, error) {
	podsToWatchCh := make(chan corev1.Pod)
	cancelPodTailCh := make(chan struct{})
	cancelCh := make(chan struct{})

	if o.Follow {
		o.cancelSignals.Watch(func() {
			close(cancelCh)
			close(cancelPodTailCh)
		})
	} else {
		close(cancelCh)
	}

	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return podsToWatchCh, cancelPodTailCh, err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return podsToWatchCh, cancelPodTailCh, err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return podsToWatchCh, cancelPodTailCh, err
	}

	watcher := NewServicePodWatcher(service, servingClient, coreClient, o.ui)

	go func() {
		watcher.Watch(podsToWatchCh, cancelCh)
		close(podsToWatchCh)
	}()

	return podsToWatchCh, cancelPodTailCh, nil
}
