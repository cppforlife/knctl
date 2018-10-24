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

package revision

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

type ShowOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	RevisionFlags cmdflags.RevisionFlags
}

func NewShowOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *ShowOptions {
	return &ShowOptions{ui: ui, depsFactory: depsFactory}
}

func NewShowCmd(o *ShowOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show revision",
		Long:  "Show revision details in a namespace",
		Example: `
  # Show details for revison 'rev1' in namespace 'ns1'
  knctl revision show -r rev1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RevisionFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *ShowOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewReference(o.RevisionFlags, tags, servingClient).Revision()
	if err != nil {
		return err
	}

	o.printStatus(revision, tags)
	o.printConditions(revision)

	podsToWatchCh, err := o.setUpPodWatching(revision)
	if err != nil {
		return err
	}

	NewPodConditionsTable(podsToWatchCh).Print(o.ui)

	return nil
}

func (o *ShowOptions) printStatus(revision *v1alpha1.Revision, tags ctlservice.Tags) {
	table := uitable.Table{
		Title: fmt.Sprintf("Revision '%s'", o.RevisionFlags.Name),

		// TODO Content: "revision",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Tags"),
			uitable.NewHeader("Serving State"),
			uitable.NewHeader("Annotations"),
			uitable.NewHeader("Age"),
		},

		Transpose: true,
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(revision.Name),
		uitable.NewValueStrings(tags.List(*revision)),
		uitable.NewValueString(string(revision.Spec.ServingState)),
		cmdcore.NewAnnotationsValue(revision.Annotations),
		cmdcore.NewValueAge(revision.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)
}

func (o *ShowOptions) printConditions(revision *v1alpha1.Revision) {
	table := uitable.Table{
		Title: fmt.Sprintf("Revision '%s' conditions", o.RevisionFlags.Name),

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

	for _, cond := range revision.Status.Conditions {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(string(cond.Type)),
			uitable.ValueFmt{
				V:     uitable.NewValueString(string(cond.Status)),
				Error: cond.Status != corev1.ConditionTrue,
			},
			cmdcore.NewValueAge(cond.LastTransitionTime.Time),
			uitable.NewValueString(cond.Reason),
			uitable.NewValueString(wordwrap.WrapString(cond.Message, 80)),
		})
	}

	o.ui.PrintTable(table)
}

func (o *ShowOptions) setUpPodWatching(revision *v1alpha1.Revision) (chan corev1.Pod, error) {
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

	watcher := ctlservice.NewRevisionPodWatcher(revision, servingClient, coreClient, o.ui)

	go func() {
		watcher.Watch(podsToWatchCh, cancelCh)
		close(podsToWatchCh)
	}()

	return podsToWatchCh, nil
}
