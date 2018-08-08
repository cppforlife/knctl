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
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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

func NewLogsCmd(o *LogsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Print logs",
		Long:  "Print logs of all active pods for a service",
		Example: `
  # Fetch last 10 log lines for service 'svc1' in namespace 'ns1' 
  knctl logs -s svc1 -n ns1

  # Follow logs for service 'svc1' in namespace 'ns1' 
  knctl logs -f -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd)

	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", false, "As new revisions are added, new pod logs will be printed")
	cmd.Flags().Int64VarP(&o.Lines, "lines", "l", 10, "Number of lines")

	return cmd
}

func (o *LogsOptions) Run() error {
	if !o.Follow && o.Lines <= 0 {
		return fmt.Errorf("Expected --lines to be greater than zero since --follow is not specified")
	}

	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	revisionWatcher := ctlservice.NewRevisionWatcher(
		servingClient.ServingV1alpha1().Revisions(o.ServiceFlags.NamespaceFlags.Name),
		metav1.ListOptions{
			LabelSelector: labels.Set(map[string]string{
				serving.ConfigurationLabelKey: o.ServiceFlags.Name,
			}).String(),
		},
	)

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	podsClient := coreClient.CoreV1().Pods(o.ServiceFlags.NamespaceFlags.Name)

	podWatcherFunc := func(revision v1alpha1.Revision) ctlservice.PodWatcher {
		return ctlservice.NewPodWatcher(
			podsClient,
			metav1.ListOptions{
				LabelSelector: labels.Set(map[string]string{
					serving.RevisionUID: string(revision.UID),
				}).String(),
			},
		)
	}

	return o.tail(revisionWatcher, podWatcherFunc, podsClient)
}

func (o *LogsOptions) tail(
	revisionWatcher ctlservice.RevisionWatcher,
	podWatcherFunc func(rev v1alpha1.Revision) ctlservice.PodWatcher,
	podsClient typedcorev1.PodInterface) error {

	cancelResWatchCh := make(chan struct{})
	cancelPodTailCh := make(chan struct{})
	doneCh := make(chan struct{})
	revisionsToWatchCh := make(chan v1alpha1.Revision)
	podsToWatchCh := make(chan corev1.Pod)

	// Watch revisions in this service
	go func() {
		err := revisionWatcher.Watch(revisionsToWatchCh, cancelResWatchCh)
		if err != nil {
			o.ui.BeginLinef("Revision watching error: %s\n", err)
		}
		close(revisionsToWatchCh)
	}()

	// Watch pods in each revision
	go func() {
		var wg sync.WaitGroup

		watchedRevs := map[string]struct{}{}

		for revision := range revisionsToWatchCh {
			revision := revision

			revUID := string(revision.UID)
			if _, found := watchedRevs[revUID]; found {
				continue
			}

			watchedRevs[revUID] = struct{}{}
			wg.Add(1)

			go func() {
				err := podWatcherFunc(revision).Watch(podsToWatchCh, cancelResWatchCh)
				if err != nil {
					o.ui.BeginLinef("Pod watching error: %s\n", err)
				}
				wg.Done()
			}()
		}

		wg.Wait()
		close(podsToWatchCh)
	}()

	// Tail logs for each pod
	go func() {
		var wg sync.WaitGroup

		watchedPods := map[string]struct{}{}
		tailOpts := logs.PodLogOpts{Follow: o.Follow}

		if o.Lines != 0 {
			tailOpts.Lines = &o.Lines
		}

		for pod := range podsToWatchCh {
			pod := pod

			podUID := string(pod.UID)
			if _, found := watchedPods[podUID]; found {
				continue
			}

			watchedPods[podUID] = struct{}{}
			wg.Add(1)

			go func() {
				tag := fmt.Sprintf("%s > %s", pod.Labels[serving.RevisionLabelKey], pod.Name)
				err := logs.NewPodContainerLog(pod, "user-container", podsClient, tag, tailOpts).Tail(o.ui, cancelPodTailCh)
				if err != nil {
					o.ui.BeginLinef("Pod logs tailing error: %s\n", err)
				}
				wg.Done()
			}()
		}

		wg.Wait()

		doneCh <- struct{}{}
	}()

	if o.Follow {
		o.cancelSignals.Watch(func() {
			close(cancelResWatchCh)
			close(cancelPodTailCh)
		})
	} else {
		close(cancelResWatchCh)
	}

	<-doneCh

	return nil
}
