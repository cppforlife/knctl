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
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cppforlife/go-cli-ui/ui"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
)

func TestNewKnctlCmd_Ok(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	realCmd := NewKnctlOptions(noopUI, NewConfigFactoryImpl(), newDepsFactory())
	cobraCmd := NewKnctlCmd(realCmd, FlagsFactory{})

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
	realCmd := NewKnctlOptions(noopUI, NewConfigFactoryImpl(), newDepsFactory())
	cmd := NewTestCmd(t, NewKnctlCmd(realCmd, FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.UIFlags, UIFlags{})
}

func TestNewKnctlCmd_OkUIFlags(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	realCmd := NewKnctlOptions(noopUI, NewConfigFactoryImpl(), newDepsFactory())
	cmd := NewTestCmd(t, NewKnctlCmd(realCmd, FlagsFactory{}))
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
	rootCmd := NewDefaultKnctlCmd(noopUI)

	const trailingSlash = " \\"

	cobrautil.VisitCommands(rootCmd, func(cmd *cobra.Command) {
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
			rootCmd := NewDefaultKnctlCmd(noopUI)

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
	rootCmd := NewDefaultKnctlCmd(noopUI)

	cobrautil.VisitCommands(rootCmd, func(cmd *cobra.Command) {
		testCmd := NewTestCmd(t, cmd)
		testCmd.ExpectBasicConfig()
	})
}

func TestNewKnctlCmd_ValidateAllCommandArgs(t *testing.T) {
	noopUI := ui.NewWrappingConfUI(ui.NewNoopUI(), ui.NewNoopLogger())
	rootCmd := NewDefaultKnctlCmd(noopUI)

	cobrautil.VisitCommands(rootCmd, func(cmd *cobra.Command) {
		if cmd.Args == nil {
			t.Fatalf("Expected command '%v' to specify arg configuration", cmd)
		}
	})
}

func TestNewKnctlCmd_ValidateAllDocsCommandExamples(t *testing.T) {
	const beginningDollar = "$ "
	const trailingSlash = " \\"

	verifyDocFile := func(content, location string) {
		lines := strings.Split(content, "\n")

		var cmdPieces []string
		var addNamespace bool

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "$ export KNCTL_NAMESPACE=") {
				addNamespace = true
				continue
			}

			if !strings.HasPrefix(line, beginningDollar+"knctl") && len(cmdPieces) == 0 {
				continue
			}

			line = strings.TrimPrefix(line, beginningDollar)

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
			rootCmd := NewDefaultKnctlCmd(noopUI)

			if cmdPieces[0] != "knctl" {
				t.Fatalf("Expected example command '%s' to start with 'knctl' (location: %s)", line, location)
			}

			if addNamespace {
				cmdPieces = append(cmdPieces, []string{"-n", "ns1"}...)
			}

			testCmd := NewTestCmd(t, rootCmd)
			testCmd.Execute(cmdPieces[1:])
			testCmd.ExpectReachesExecution()

			cmdPieces = []string{}
		}
	}

	matches, err := filepath.Glob("../../../docs/*.md")
	if err != nil {
		t.Fatalf("Expected glob to not error: %s", err)
	}

	if len(matches) == 0 {
		t.Fatalf("Expected glob to find at least one doc file")
	}

	for _, match := range matches {
		content, err := ioutil.ReadFile(match)
		if err != nil {
			t.Fatalf("Expected file reading to succeed: %s", err)
		}

		verifyDocFile(string(content), match)
	}
}
