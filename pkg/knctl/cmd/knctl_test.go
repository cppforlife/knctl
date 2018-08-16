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

package cmd_test

import (
	"strings"
	"testing"

	"github.com/cppforlife/go-cli-ui/ui"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	"github.com/spf13/cobra"
)

func TestNewKnctlCmd_Ok(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	realCmd := NewDefaultKnctlOptions(noopUI)
	cobraCmd := NewKnctlCmd(realCmd)

	cmd := NewTestCmd(t, cobraCmd)
	cmd.ExpectBasicConfig()

	if !cobraCmd.SilenceErrors {
		t.Fatalf("Expected SilenceErrors to be true")
	}
	if !cobraCmd.SilenceUsage {
		t.Fatalf("Expected SilenceUsage to be true")
	}
}

func TestNewKnctlCmd_OkMinimum(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	realCmd := NewDefaultKnctlOptions(noopUI)
	cmd := NewTestCmd(t, NewKnctlCmd(realCmd))
	cmd.Execute([]string{})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.UIFlags, UIFlags{})
}

func TestNewKnctlCmd_OkUIFlags(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	realCmd := NewDefaultKnctlOptions(noopUI)
	cmd := NewTestCmd(t, NewKnctlCmd(realCmd))
	cmd.Execute([]string{
		"--tty",
		"--no-color",
		"--json",
		"--non-interactive",
		"--column", "test-col1",
		"--column", "test-col2",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.UIFlags, UIFlags{
		TTY:            true,
		NoColor:        true,
		JSON:           true,
		NonInteractive: true,
		Columns:        []string{"test-col1", "test-col2"},
	})
}

func TestNewKnctlCmd_ValidateAllCommandExamples(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	rootCmd := NewKnctlCmd(NewDefaultKnctlOptions(noopUI))

	const trailingSlash = " \\"

	VisitCommands(rootCmd, func(cmd *cobra.Command) {
		lines := strings.Split(cmd.Example, "\n")

		var cmdPieces []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 || strings.HasPrefix(line, "#") {
				continue
			}

			var endsWithSlash bool

			if strings.HasSuffix(line, trailingSlash) {
				line = strings.TrimSuffix(line, trailingSlash)
				endsWithSlash = true
			}

			cmdPieces = append(cmdPieces, strings.Split(line, " ")...)
			if endsWithSlash {
				continue
			}

			// recreate for every command since cobra persists some state
			noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
			rootCmd := NewKnctlCmd(NewDefaultKnctlOptions(noopUI))

			if cmdPieces[0] != "knctl" {
				t.Fatalf("Expected example command '%s' to start with 'knctl'", line)
			}

			testCmd := NewTestCmd(t, rootCmd)
			testCmd.Execute(cmdPieces[1:])
			testCmd.ExpectReachesExecution()

			cmdPieces = []string{}
		}
	})
}

func TestNewKnctlCmd_ValidateAllCommandBasicConfig(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	rootCmd := NewKnctlCmd(NewDefaultKnctlOptions(noopUI))

	VisitCommands(rootCmd, func(cmd *cobra.Command) {
		testCmd := NewTestCmd(t, cmd)
		testCmd.ExpectBasicConfig()
	})
}

func TestNewKnctlCmd_ValidateAllCommandArgs(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	rootCmd := NewKnctlCmd(NewDefaultKnctlOptions(noopUI))

	VisitCommands(rootCmd, func(cmd *cobra.Command) {
		if cmd.Args == nil {
			t.Fatalf("Expected command '%v' to specify arg configuration", cmd)
		}
	})
}
