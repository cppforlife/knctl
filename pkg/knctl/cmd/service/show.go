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

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	cmdrev "github.com/cppforlife/knctl/pkg/knctl/cmd/revision"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShowOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ServiceFlags cmdflags.ServiceFlags
}

func NewShowOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *ShowOptions {
	return &ShowOptions{ui: ui, depsFactory: depsFactory}
}

func NewShowCmd(o *ShowOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
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

func (o *ShowOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	o.printStatus(service)

	cmdcore.NewConditionsTable(service.Status.Conditions).Print(o.ui)

	podsToWatchCh, err := o.setUpPodWatching(service)
	if err != nil {
		return err
	}

	cmdrev.NewPodConditionsTable(podsToWatchCh).Print(o.ui)

	return nil
}

func (o ShowOptions) printStatus(service *v1alpha1.Service) {
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
		cmdcore.NewAnnotationsValue(service.Annotations),
		cmdcore.NewValueAge(service.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)
}

func (o *ShowOptions) setUpPodWatching(service *v1alpha1.Service) (chan corev1.Pod, error) {
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

	watcher := ctlservice.NewServicePodWatcher(service, servingClient, coreClient, o.ui)

	go func() {
		watcher.Watch(podsToWatchCh, cancelCh)
		close(podsToWatchCh)
	}()

	return podsToWatchCh, nil
}
