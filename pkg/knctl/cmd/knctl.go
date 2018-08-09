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
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
)

type KnctlOptions struct {
	ui          *ui.ConfUI
	depsFactory DepsFactory

	UIFlags         UIFlags
	KubeconfigFlags KubeconfigFlags
}

func NewDefaultKnctlOptions(ui *ui.ConfUI) *KnctlOptions {
	return NewKnctlOptions(ui, NewDepsFactoryImpl())
}

func NewKnctlOptions(ui *ui.ConfUI, depsFactory DepsFactory) *KnctlOptions {
	return &KnctlOptions{ui: ui, depsFactory: depsFactory}
}

func NewKnctlCmd(o *KnctlOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knctl",
		Short: "knctl controls Knative resources",
		Long:  "Find more information at: https://github.com/knative/docs",

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			o.depsFactory.ConfigureConfigPath(o.KubeconfigFlags.Path)
			o.UIFlags.ConfigureUI(o.ui)
		},

		RunE: ShowHelp,

		// Affects children as well
		SilenceErrors: true,
		SilenceUsage:  true,

		// Disable docs header
		DisableAutoGenTag: true,

		// TODO bash completion
	}

	o.UIFlags.Set(cmd)
	o.KubeconfigFlags.Set(cmd)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui)))
	cmd.AddCommand(NewInstallCmd(NewInstallOptions(o.ui, o.depsFactory, &o.KubeconfigFlags)))
	cmd.AddCommand(NewDeployCmd(NewDeployOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(NewLogsCmd(NewLogsOptions(o.ui, o.depsFactory, CancelSignals{})))
	cmd.AddCommand(NewOpenCmd(NewOpenOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(NewCurlCmd(NewCurlOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(NewBuildCmd(NewBuildOptions(o.ui, o.depsFactory, CancelSignals{})))
	cmd.AddCommand(NewRouteCmd(NewRouteOptions(o.ui, o.depsFactory)))

	createCmd := NewCreateCmd()
	createCmd.AddCommand(NewCreateNamespaceCmd(NewCreateNamespaceOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(createCmd)

	listCmd := NewListCmd()
	listCmd.AddCommand(NewListServicesCmd(NewListServicesOptions(o.ui, o.depsFactory)))
	listCmd.AddCommand(NewListRevisionsCmd(NewListRevisionsOptions(o.ui, o.depsFactory)))
	listCmd.AddCommand(NewListPodsCmd(NewListPodsOptions(o.ui, o.depsFactory)))
	listCmd.AddCommand(NewListBuildsCmd(NewListBuildsOptions(o.ui, o.depsFactory)))
	listCmd.AddCommand(NewListRoutesCmd(NewListRoutesOptions(o.ui, o.depsFactory)))
	listCmd.AddCommand(NewListIngressesCmd(NewListIngressesOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(listCmd)

	deleteCmd := NewDeleteCmd()
	deleteCmd.AddCommand(NewDeleteServiceCmd(NewDeleteServiceOptions(o.ui, o.depsFactory)))
	deleteCmd.AddCommand(NewDeleteRevisionCmd(NewDeleteRevisionOptions(o.ui, o.depsFactory)))
	deleteCmd.AddCommand(NewDeleteBuildCmd(NewDeleteBuildOptions(o.ui, o.depsFactory)))
	deleteCmd.AddCommand(NewDeleteRouteCmd(NewDeleteRouteOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(deleteCmd)

	annotateCmd := NewAnnotateCmd()
	annotateCmd.AddCommand(NewAnnotateServiceCmd(NewAnnotateServiceOptions(o.ui, o.depsFactory)))
	annotateCmd.AddCommand(NewAnnotateRevisionCmd(NewAnnotateRevisionOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(annotateCmd)

	tagCmd := NewTagCmd()
	tagCmd.AddCommand(NewTagRevisionCmd(NewTagRevisionOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(tagCmd)

	untagCmd := NewUntagCmd()
	untagCmd.AddCommand(NewUntagRevisionCmd(NewUntagRevisionOptions(o.ui, o.depsFactory)))
	cmd.AddCommand(untagCmd)

	VisitCommands(cmd, func(c *cobra.Command) {
		if c.Args == nil {
			c.Args = cobra.NoArgs
		}
	})

	return cmd
}

func ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.Help() // TODO use ui for output
	return nil
}

func VisitCommands(cmd *cobra.Command, f func(*cobra.Command)) {
	f(cmd)
	for _, child := range cmd.Commands() {
		VisitCommands(child, f)
	}
}
