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
	"reflect"
	"strings"
	"testing"

	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
)

type TestCmd struct {
	t   *testing.T
	cmd *cobra.Command

	executeCalled         bool
	executeArgs           []string
	executeErr            error
	didReachNoopExecution bool // noop execution reached
}

func NewTestCmd(t *testing.T, cmd *cobra.Command) *TestCmd {
	return &TestCmd{t: t, cmd: cmd}
}

func (c *TestCmd) ExpectBasicConfig() {
	if c.executeCalled {
		c.t.Fatalf("Expected Execute() to not be called before checking basic config")
	}

	if len(c.cmd.Use) == 0 {
		c.t.Fatalf("Expected command to have 'Use' set")
	}
	if len(c.cmd.Short) == 0 {
		c.t.Fatalf("Expected command to have 'Short' set")
	}
	if c.cmd.RunE == nil {
		c.t.Fatalf("Expected command to have 'RunE' set")
	}
}

func (c *TestCmd) Execute(args []string) {
	if c.executeCalled {
		c.t.Fatalf("Expected Execute() to not be called multiple times")
	}

	c.executeCalled = true

	cobrautil.VisitCommands(c.cmd, func(ci *cobra.Command) {
		ci.RunE = func(_ *cobra.Command, _ []string) error {
			c.didReachNoopExecution = true
			return nil
		}
	})

	c.cmd.SilenceErrors = true
	c.cmd.SilenceUsage = true

	c.cmd.SetArgs(args)
	c.executeArgs = args
	c.executeErr = c.cmd.Execute()
}

func (c *TestCmd) expectExecuteCalled() {
	if !c.executeCalled {
		c.t.Fatalf("Expected Execute() to be called before assertion")
	}
}

func (c *TestCmd) ExpectRequiredFlags(flags []string) {
	c.expectExecuteCalled()

	if len(flags) == 0 {
		c.t.Fatalf("Expected at least one required flag")
	}

	if c.executeErr == nil {
		c.t.Fatalf("Expected execute error")
	}
	if c.didReachNoopExecution {
		c.t.Fatalf("Expected command to not reach execution")
	}

	expectedErrMsg := fmt.Sprintf(`required flag(s) "%s" not set`, strings.Join(flags, `", "`))

	if c.executeErr.Error() != expectedErrMsg {
		c.t.Fatalf("Expected required flags error ('%s'), but was '%s'", expectedErrMsg, c.executeErr)
	}
}

func (c *TestCmd) ExpectReachesExecution() {
	c.expectExecuteCalled()

	if c.executeErr != nil {
		c.t.Fatalf("[command '%v'] Expected nil error, but was: '%s'", c.executeArgs, c.executeErr)
	}
	if !c.didReachNoopExecution {
		c.t.Fatalf("[command '%v'] Expected command to reach execution", c.executeArgs)
	}
}

func DeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expect obj '%#v' to equal obj '%#v'", actual, expected)
	}
}
