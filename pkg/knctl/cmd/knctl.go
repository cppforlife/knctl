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
	"io"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
)

type KnctlOptions struct {
	ui            *ui.ConfUI
	configFactory ConfigFactory
	depsFactory   DepsFactory

	UIFlags         UIFlags
	KubeconfigFlags KubeconfigFlags
}

func NewKnctlOptions(ui *ui.ConfUI, configFactory ConfigFactory, depsFactory DepsFactory) *KnctlOptions {
	return &KnctlOptions{ui: ui, configFactory: configFactory, depsFactory: depsFactory}
}

func NewDefaultKnctlCmd(ui *ui.ConfUI) *cobra.Command {
	configFactory := NewConfigFactoryImpl()
	depsFactory := NewDepsFactoryImpl(configFactory)
	options := NewKnctlOptions(ui, configFactory, depsFactory)
	flagsFactory := NewFlagsFactory(configFactory, depsFactory)
	return NewKnctlCmd(options, flagsFactory)
}

func NewKnctlCmd(o *KnctlOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knctl",
		Short: "knctl controls Knative resources",
		Long: `knctl controls Knative resources.

CLI docs: https://github.com/cppforlife/knctl#docs.
Knative docs: https://github.com/knative/docs.`,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			o.configFactory.ConfigurePath(o.KubeconfigFlags.Path)
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

	cmd.SetOutput(uiBlockWriter{o.ui}) // setting output for cmd.Help()

	o.UIFlags.Set(cmd, flagsFactory)
	o.KubeconfigFlags.Set(cmd, flagsFactory)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui), flagsFactory))
	cmd.AddCommand(NewInstallCmd(NewInstallOptions(o.ui, o.depsFactory, &o.KubeconfigFlags), flagsFactory))
	cmd.AddCommand(NewUninstallCmd(NewUninstallOptions(o.ui, o.depsFactory, &o.KubeconfigFlags), flagsFactory))
	cmd.AddCommand(NewDeployCmd(NewDeployOptions(o.ui, o.configFactory, o.depsFactory), flagsFactory))
	cmd.AddCommand(NewLogsCmd(NewLogsOptions(o.ui, o.depsFactory, CancelSignals{}), flagsFactory))
	cmd.AddCommand(NewOpenCmd(NewOpenOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(NewCurlCmd(NewCurlOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(NewBuildCmd(NewBuildOptions(o.ui, o.configFactory, o.depsFactory, CancelSignals{}), flagsFactory))
	cmd.AddCommand(NewRouteCmd(NewRouteOptions(o.ui, o.depsFactory), flagsFactory))

	createCmd := NewCreateCmd()
	createCmd.AddCommand(NewCreateNamespaceCmd(NewCreateNamespaceOptions(o.ui, o.depsFactory), flagsFactory))
	createCmd.AddCommand(NewCreateServiceAccountCmd(NewCreateServiceAccountOptions(o.ui, o.depsFactory), flagsFactory))
	createCmd.AddCommand(NewCreateBasicAuthSecretCmd(NewCreateBasicAuthSecretOptions(o.ui, o.depsFactory), flagsFactory))
	createCmd.AddCommand(NewCreateSSHAuthSecretCmd(NewCreateSSHAuthSecretOptions(o.ui, o.depsFactory), flagsFactory))
	createCmd.AddCommand(NewCreateDomainCmd(NewCreateDomainOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(createCmd)

	listCmd := NewListCmd()
	listCmd.AddCommand(NewListServicesCmd(NewListServicesOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListRevisionsCmd(NewListRevisionsOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListPodsCmd(NewListPodsOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListBuildsCmd(NewListBuildsOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListRoutesCmd(NewListRoutesOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListDomainsCmd(NewListDomainsOptions(o.ui, o.depsFactory), flagsFactory))
	listCmd.AddCommand(NewListIngressesCmd(NewListIngressesOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(listCmd)

	deleteCmd := NewDeleteCmd()
	deleteCmd.AddCommand(NewDeleteServiceCmd(NewDeleteServiceOptions(o.ui, o.depsFactory), flagsFactory))
	deleteCmd.AddCommand(NewDeleteRevisionCmd(NewDeleteRevisionOptions(o.ui, o.depsFactory), flagsFactory))
	deleteCmd.AddCommand(NewDeleteBuildCmd(NewDeleteBuildOptions(o.ui, o.depsFactory), flagsFactory))
	deleteCmd.AddCommand(NewDeleteRouteCmd(NewDeleteRouteOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(deleteCmd)

	annotateCmd := NewAnnotateCmd()
	annotateCmd.AddCommand(NewAnnotateServiceCmd(NewAnnotateServiceOptions(o.ui, o.depsFactory), flagsFactory))
	annotateCmd.AddCommand(NewAnnotateRevisionCmd(NewAnnotateRevisionOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(annotateCmd)

	tagCmd := NewTagCmd()
	tagCmd.AddCommand(NewTagRevisionCmd(NewTagRevisionOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(tagCmd)

	untagCmd := NewUntagCmd()
	untagCmd.AddCommand(NewUntagRevisionCmd(NewUntagRevisionOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(untagCmd)

	cobrautil.VisitCommands(cmd, reconfigureCmdWithSubcmd)
	cobrautil.VisitCommands(cmd, reconfigureLeafCmd)
	cobrautil.VisitCommands(cmd, cobrautil.ResolveFlagsForCmd)

	return cmd
}

func reconfigureCmdWithSubcmd(cmd *cobra.Command) {
	if len(cmd.Commands()) == 0 {
		return
	}

	if cmd.Args == nil {
		cmd.Args = cobra.ArbitraryArgs
	}
	if cmd.RunE == nil {
		cmd.RunE = ShowSubcommands
	}

	var strs []string
	for _, subcmd := range cmd.Commands() {
		strs = append(strs, subcmd.Use)
	}

	cmd.Short += " (" + strings.Join(strs, ", ") + ")"
}

func reconfigureLeafCmd(cmd *cobra.Command) {
	if len(cmd.Commands()) > 0 {
		return
	}

	if cmd.RunE == nil {
		panic(fmt.Sprintf("Internal: Command '%s' does not set RunE", cmd.CommandPath()))
	}

	if cmd.Args == nil {
		origRunE := cmd.RunE
		cmd.RunE = func(cmd2 *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("command '%s' does not accept extra arguments '%s'", args[0], cmd2.CommandPath())
			}
			return origRunE(cmd2, args)
		}
		cmd.Args = cobra.ArbitraryArgs
	}
}

func ShowSubcommands(cmd *cobra.Command, args []string) error {
	var strs []string
	for _, subcmd := range cmd.Commands() {
		strs = append(strs, subcmd.Use)
	}
	return fmt.Errorf("Use one of available subcommands: %s", strings.Join(strs, ", "))
}

func ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return fmt.Errorf("Invalid command - see available commands/subcommands above")
}

type uiBlockWriter struct {
	ui ui.UI
}

var _ io.Writer = uiBlockWriter{}

func (w uiBlockWriter) Write(p []byte) (n int, err error) {
	w.ui.PrintBlock(p)
	return len(p), nil
}
