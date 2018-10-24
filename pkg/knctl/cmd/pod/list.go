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

package pod

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ServiceFlags cmdflags.ServiceFlags
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: cmdcore.ListAliases,
		Short:   "List pods",
		Long:    "List all pods for a service",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *ListOptions) Run() error {
	podsToWatchCh, err := o.setUpPodWatching()
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: fmt.Sprintf("Pods for service '%s'", o.ServiceFlags.Name),

		Content: "pods",

		Header: []uitable.Header{
			uitable.NewHeader("Revision"),
			uitable.NewHeader("Name"),
			uitable.NewHeader("Phase"),
			uitable.NewHeader("Restarts"),
			uitable.NewHeader("Age"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for pod := range podsToWatchCh {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(pod.Labels[serving.RevisionLabelKey]),
			uitable.NewValueString(pod.Name),
			uitable.ValueFmt{
				V:     uitable.NewValueString(string(pod.Status.Phase)),
				Error: !o.isOKStatus(pod),
			},
			uitable.NewValueInt(o.podRestarts(pod)),
			cmdcore.NewValueAge(pod.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *ListOptions) isOKStatus(pod corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodSucceeded
}

func (o *ListOptions) podRestarts(pod corev1.Pod) int {
	var count int
	for _, status := range pod.Status.ContainerStatuses {
		count += int(status.RestartCount)
	}
	return count
}

func (o *ListOptions) setUpPodWatching() (chan corev1.Pod, error) {
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

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
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
