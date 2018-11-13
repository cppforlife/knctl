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

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShowOptions struct {
	ui            ui.UI
	configFactory cmdcore.ConfigFactory
	depsFactory   cmdcore.DepsFactory

	BuildFlags BuildFlags
	Logs       bool
}

func NewShowOptions(ui ui.UI, configFactory cmdcore.ConfigFactory, depsFactory cmdcore.DepsFactory) *ShowOptions {
	return &ShowOptions{ui: ui, configFactory: configFactory, depsFactory: depsFactory}
}

func NewShowCmd(o *ShowOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
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

func (o *ShowOptions) Run() error {
	buildClient, err := o.depsFactory.BuildClient()
	if err != nil {
		return err
	}

	build, err := buildClient.BuildV1alpha1().Builds(o.BuildFlags.NamespaceFlags.Name).Get(o.BuildFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	o.printStatus(build)

	cmdcore.NewConditionsTable(build.Status.Conditions).Print(o.ui)

	if o.Logs {
		return o.showLogs(build, buildClient)
	}

	return nil
}

func (o *ShowOptions) printStatus(build *v1alpha1.Build) {
	table := uitable.Table{
		Title: fmt.Sprintf("Build '%s'", o.BuildFlags.Name),
		// TODO Content: "build",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Timeout"),
			uitable.NewHeader("Started at"),
			uitable.NewHeader("Completed at"),
			uitable.NewHeader("Succeeded"),
			uitable.NewHeader("Age"),
		},

		Transpose: true,
	}

	durationVal := uitable.NewValueString("")

	if build.Spec.Timeout != nil {
		durationVal = uitable.NewValueString(build.Spec.Timeout.Duration.String())
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(build.Name),
		durationVal,
		uitable.NewValueTime(build.Status.StartTime.Time),
		uitable.NewValueTime(build.Status.CompletionTime.Time),
		NewBuildSucceededValue(*build),
		cmdcore.NewValueAge(build.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)
}

func (o *ShowOptions) showLogs(build *v1alpha1.Build, buildClient buildclientset.Interface) error {
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

	return buildObjFactory.New(build).TailLogs(o.ui, cancelCh)
}
