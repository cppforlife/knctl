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
	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShowBuildOptions struct {
	ui            ui.UI
	configFactory ConfigFactory
	depsFactory   DepsFactory
	cancelSignals CancelSignals

	BuildFlags BuildFlags
	Logs       bool
}

func NewShowBuildOptions(ui ui.UI, configFactory ConfigFactory, depsFactory DepsFactory, cancelSignals CancelSignals) *ShowBuildOptions {
	return &ShowBuildOptions{ui: ui, configFactory: configFactory, depsFactory: depsFactory, cancelSignals: cancelSignals}
}

func NewShowBuildCmd(o *ShowBuildOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show build",
		Long:  "Show build details in a namespace",
		Example: `
  # Show details for build 'build1' in namespace 'ns1'
  knctl build show -b build1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.BuildFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVar(&o.Logs, "logs", true, "Show logs")
	return cmd
}

func (o *ShowBuildOptions) Run() error {
	buildClient, err := o.depsFactory.BuildClient()
	if err != nil {
		return err
	}

	build, err := buildClient.BuildV1alpha1().Builds(o.BuildFlags.NamespaceFlags.Name).Get(o.BuildFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: fmt.Sprintf("Build '%s'", o.BuildFlags.Name),

		// TODO Content: "build",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Succeeded"),
			uitable.NewHeader("Age"),
		},

		Transpose: true,
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(build.Name),
		NewBuildSucceededValue(*build),
		NewValueAge(build.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)

	table = uitable.Table{
		Title: fmt.Sprintf("Build '%s' conditions", o.BuildFlags.Name),

		// TODO Content: "conditions",

		Header: []uitable.Header{
			uitable.NewHeader("Type"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Reason"),
			uitable.NewHeader("Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, cond := range build.Status.Conditions {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(string(cond.Type)),
			uitable.ValueFmt{
				V:     uitable.NewValueString(string(cond.Status)),
				Error: cond.Status != corev1.ConditionTrue,
			},
			// TODO age
			uitable.NewValueString(cond.Reason),
			uitable.NewValueString(wordwrap.WrapString(cond.Message, 80)),
		})
	}

	o.ui.PrintTable(table)

	if o.Logs {
		return o.showLogs(build, buildClient)
	}

	return nil
}

func (o *ShowBuildOptions) showLogs(build *v1alpha1.Build, buildClient buildclientset.Interface) error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	restConfig, err := o.configFactory.RESTConfig()
	if err != nil {
		return err
	}

	buildObjFactory := ctlbuild.NewFactory(buildClient, coreClient, restConfig)

	cancelCh := make(chan struct{})
	o.cancelSignals.Watch(func() { close(cancelCh) })

	return buildObjFactory.New(build).TailLogs(o.ui, cancelCh)
}
