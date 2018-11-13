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
	cmdbas "github.com/cppforlife/knctl/pkg/knctl/cmd/basicauthsecret"
	cmdbld "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmddom "github.com/cppforlife/knctl/pkg/knctl/cmd/domain"
	cmding "github.com/cppforlife/knctl/pkg/knctl/cmd/ingress"
	cmdkn "github.com/cppforlife/knctl/pkg/knctl/cmd/knative"
	cmdpod "github.com/cppforlife/knctl/pkg/knctl/cmd/pod"
	cmdrev "github.com/cppforlife/knctl/pkg/knctl/cmd/revision"
	cmdrte "github.com/cppforlife/knctl/pkg/knctl/cmd/route"
	cmdsvc "github.com/cppforlife/knctl/pkg/knctl/cmd/service"
	cmdsa "github.com/cppforlife/knctl/pkg/knctl/cmd/serviceaccount"
	cmdsas "github.com/cppforlife/knctl/pkg/knctl/cmd/sshauthsecret"
	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
)

type KnctlOptions struct {
	ui            *ui.ConfUI
	configFactory cmdcore.ConfigFactory
	depsFactory   cmdcore.DepsFactory

	UIFlags         cmdcore.UIFlags
	KubeconfigFlags cmdcore.KubeconfigFlags
}

func NewKnctlOptions(ui *ui.ConfUI, configFactory cmdcore.ConfigFactory, depsFactory cmdcore.DepsFactory) *KnctlOptions {
	return &KnctlOptions{ui: ui, configFactory: configFactory, depsFactory: depsFactory}
}

func NewDefaultKnctlCmd(ui *ui.ConfUI) *cobra.Command {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory)
	options := NewKnctlOptions(ui, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)
	return NewKnctlCmd(options, flagsFactory)
}

func NewKnctlCmd(o *KnctlOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knctl",
		Short: "knctl controls Knative resources",
		Long: `knctl controls Knative resources.

CLI docs: https://github.com/cppforlife/knctl#docs.
Knative docs: https://github.com/knative/docs.`,

		RunE: ShowHelp,

		// Affects children as well
		SilenceErrors: true,
		SilenceUsage:  true,

		// Disable docs header
		DisableAutoGenTag: true,

		// TODO bash completion
	}

	cmd.SetOutput(uiBlockWriter{o.ui}) // setting output for cmd.Help()

	cmd.SetUsageTemplate(cobrautil.HelpSectionsUsageTemplate([]cobrautil.HelpSection{
		cmdcore.BasicHelpGroup,
		cmdcore.BuildMgmtHelpGroup,
		cmdcore.SecretMgmtHelpGroup,
		cmdcore.RouteMgmtHelpGroup,
		cmdcore.OtherHelpGroup,
		cmdcore.SystemHelpGroup,
		cmdcore.RestOfCommandsHelpGroup,
	}))

	o.UIFlags.Set(cmd, flagsFactory)
	o.KubeconfigFlags.Set(cmd, flagsFactory)

	o.configFactory.ConfigurePathResolver(o.KubeconfigFlags.Path.Value)
	o.configFactory.ConfigureContextResolver(o.KubeconfigFlags.Context.Value)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui), flagsFactory))

	// Knative
	cmd.AddCommand(cmdkn.NewInstallCmd(cmdkn.NewInstallOptions(o.ui, o.depsFactory, &o.KubeconfigFlags), flagsFactory))
	cmd.AddCommand(cmdkn.NewUninstallCmd(cmdkn.NewUninstallOptions(o.ui, o.depsFactory), flagsFactory))

	serviceCmd := cmdsvc.NewCmd()
	serviceCmd.AddCommand(cmdsvc.NewListCmd(cmdsvc.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	serviceCmd.AddCommand(cmdsvc.NewShowCmd(cmdsvc.NewShowOptions(o.ui, o.depsFactory), flagsFactory))
	serviceCmd.AddCommand(cmdsvc.NewDeleteCmd(cmdsvc.NewDeleteOptions(o.ui, o.depsFactory), flagsFactory))
	serviceCmd.AddCommand(cmdsvc.NewAnnotateCmd(cmdsvc.NewAnnotateOptions(o.ui, o.depsFactory), flagsFactory))
	serviceCmd.AddCommand(cmdsvc.NewOpenCmd(cmdsvc.NewOpenOptions(o.ui, o.depsFactory), flagsFactory))
	serviceCmd.AddCommand(cmdsvc.NewURLCmd(cmdsvc.NewURLOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(serviceCmd)

	cmd.AddCommand(cmdsvc.NewDeployCmd(cmdsvc.NewDeployOptions(o.ui, o.configFactory, o.depsFactory), flagsFactory))
	cmd.AddCommand(cmdsvc.NewLogsCmd(cmdsvc.NewLogsOptions(o.ui, o.depsFactory, cmdcore.CancelSignals{}), flagsFactory))
	cmd.AddCommand(cmdsvc.NewCurlCmd(cmdsvc.NewCurlOptions(o.ui, o.depsFactory), flagsFactory))

	revisionCmd := cmdrev.NewCmd()
	revisionCmd.AddCommand(cmdrev.NewListCmd(cmdrev.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	revisionCmd.AddCommand(cmdrev.NewShowCmd(cmdrev.NewShowOptions(o.ui, o.depsFactory), flagsFactory))
	revisionCmd.AddCommand(cmdrev.NewDeleteCmd(cmdrev.NewDeleteOptions(o.ui, o.depsFactory), flagsFactory))
	revisionCmd.AddCommand(cmdrev.NewTagCmd(cmdrev.NewTagOptions(o.ui, o.depsFactory), flagsFactory))
	revisionCmd.AddCommand(cmdrev.NewUntagCmd(cmdrev.NewUntagOptions(o.ui, o.depsFactory), flagsFactory))
	revisionCmd.AddCommand(cmdrev.NewAnnotateCmd(cmdrev.NewAnnotateOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(revisionCmd)

	routeCmd := cmdrte.NewCmd()
	routeCmd.AddCommand(cmdrte.NewCreateCmd(cmdrte.NewCreateOptions(o.ui, o.depsFactory), flagsFactory))
	routeCmd.AddCommand(cmdrte.NewShowCmd(cmdrte.NewShowOptions(o.ui, o.depsFactory), flagsFactory))
	routeCmd.AddCommand(cmdrte.NewListCmd(cmdrte.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	routeCmd.AddCommand(cmdrte.NewDeleteCmd(cmdrte.NewDeleteOptions(o.ui, o.depsFactory), flagsFactory))
	routeCmd.AddCommand(cmdrte.NewCurlCmd(cmdrte.NewCurlOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(routeCmd)

	buildCmd := cmdbld.NewCmd()
	buildCmd.AddCommand(cmdbld.NewCreateCmd(cmdbld.NewCreateOptions(o.ui, o.configFactory, o.depsFactory), flagsFactory))
	buildCmd.AddCommand(cmdbld.NewListCmd(cmdbld.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	buildCmd.AddCommand(cmdbld.NewShowCmd(cmdbld.NewShowOptions(o.ui, o.configFactory, o.depsFactory), flagsFactory))
	buildCmd.AddCommand(cmdbld.NewDeleteCmd(cmdbld.NewDeleteOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(buildCmd)

	domainCmd := cmddom.NewCmd()
	domainCmd.AddCommand(cmddom.NewCreateCmd(cmddom.NewCreateOptions(o.ui, o.depsFactory), flagsFactory))
	domainCmd.AddCommand(cmddom.NewListCmd(cmddom.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(domainCmd)

	cmd.AddCommand(cmddom.NewDNSMapCmd(cmddom.NewDNSMapOptions(o.ui, o.depsFactory), flagsFactory))

	ingressCmd := cmding.NewCmd()
	ingressCmd.AddCommand(cmding.NewListCmd(cmding.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(ingressCmd)

	podCmd := cmdpod.NewCmd()
	podCmd.AddCommand(cmdpod.NewListCmd(cmdpod.NewListOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(podCmd)

	serviceAccountCmd := cmdsa.NewCmd()
	serviceAccountCmd.AddCommand(cmdsa.NewCreateCmd(cmdsa.NewCreateOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(serviceAccountCmd)

	basicAuthSecretCmd := cmdbas.NewCmd()
	basicAuthSecretCmd.AddCommand(cmdbas.NewCreateCmd(cmdbas.NewCreateOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(basicAuthSecretCmd)

	sshAuthSecretCmd := cmdsas.NewCmd()
	sshAuthSecretCmd.AddCommand(cmdsas.NewCreateCmd(cmdsas.NewCreateOptions(o.ui, o.depsFactory), flagsFactory))
	cmd.AddCommand(sshAuthSecretCmd)

	// Last one runs first
	cobrautil.VisitCommands(cmd, reconfigureCmdWithSubcmd)
	cobrautil.VisitCommands(cmd, reconfigureLeafCmd)

	cobrautil.VisitCommands(cmd, cobrautil.WrapRunEForCmd(func(*cobra.Command, []string) error {
		o.UIFlags.ConfigureUI(o.ui)
		return nil
	}))

	cobrautil.VisitCommands(cmd, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))

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
