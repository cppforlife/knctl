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

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShowServiceOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags ServiceFlags
}

func NewShowServiceOptions(ui ui.UI, depsFactory DepsFactory) *ShowServiceOptions {
	return &ShowServiceOptions{ui: ui, depsFactory: depsFactory}
}

func NewShowServiceCmd(o *ShowServiceOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show service",
		Long:  "Show service details in a namespace",
		Example: `
  # Show details for service 'srv1' in namespace 'ns1'
  knctl service show -s srv1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *ShowServiceOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	o.printStatus(service)
	o.printConditions(service)

	podsToWatchCh, err := o.setUpPodWatching()
	if err != nil {
		return err
	}

	for pod := range podsToWatchCh {
		o.printPodConditions(pod)
	}

	return nil
}

func (o ShowServiceOptions) printStatus(service *v1alpha1.Service) {
	table := uitable.Table{
		Title: fmt.Sprintf("Service '%s'", o.ServiceFlags.Name),

		// TODO Content: "service",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Domain"),
			uitable.NewHeader("Internal Domain"),
			uitable.NewHeader("Annotations"),
			uitable.NewHeader("Age"),
		},

		Transpose: true,
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(service.Name),
		uitable.NewValueString(service.Status.Domain),
		uitable.NewValueString(service.Status.DomainInternal),
		NewAnnotationsValue(service.Annotations),
		NewValueAge(service.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)
}

func (o *ShowServiceOptions) printConditions(service *v1alpha1.Service) {
	table := uitable.Table{
		Title: fmt.Sprintf("Service '%s' conditions", o.ServiceFlags.Name),

		// TODO Content: "conditions",

		Header: []uitable.Header{
			uitable.NewHeader("Type"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Age"),
			uitable.NewHeader("Reason"),
			uitable.NewHeader("Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, cond := range service.Status.Conditions {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(string(cond.Type)),
			uitable.ValueFmt{
				V:     uitable.NewValueString(string(cond.Status)),
				Error: cond.Status != corev1.ConditionTrue,
			},
			NewValueAge(cond.LastTransitionTime.Time),
			uitable.NewValueString(cond.Reason),
			uitable.NewValueString(wordwrap.WrapString(cond.Message, 80)),
		})
	}

	o.ui.PrintTable(table)
}

func (o *ShowServiceOptions) setUpPodWatching() (chan corev1.Pod, error) {
	podsToWatchCh := make(chan corev1.Pod)
	cancelCh := make(chan struct{})
	close(cancelCh) // Close immediately for just plain listing of revisions and pods

	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return podsToWatchCh, err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return podsToWatchCh, err
	}

	watcher := NewRevisionPodWatcher(
		o.ServiceFlags.NamespaceFlags.Name, o.ServiceFlags.Name, servingClient, coreClient, o.ui)

	go func() {
		watcher.Watch(podsToWatchCh, cancelCh)
		close(podsToWatchCh)
	}()

	return podsToWatchCh, nil
}

func (o *ShowServiceOptions) printPodConditions(pod corev1.Pod) {
	table := uitable.Table{
		Title: fmt.Sprintf("Pod '%s' conditions", pod.Name),

		// TODO Content: "conditions",

		Header: []uitable.Header{
			uitable.NewHeader("Type"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Age"),
			uitable.NewHeader("Reason"),
			uitable.NewHeader("Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, cond := range pod.Status.Conditions {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(string(cond.Type)),
			uitable.ValueFmt{
				V:     uitable.NewValueString(string(cond.Status)),
				Error: cond.Status != corev1.ConditionTrue,
			},
			NewValueAge(cond.LastTransitionTime.Time),
			uitable.NewValueString(cond.Reason),
			uitable.NewValueString(wordwrap.WrapString(cond.Message, 80)),
		})
	}

	o.ui.PrintTable(table)
}
