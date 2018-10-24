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

package build_test

import (
	"testing"

	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
)

func TestNewListCmd_Ok(t *testing.T) {
	realCmd := NewListOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewListCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.NamespaceFlags, cmdcore.NamespaceFlags{"test-namespace"})
}

func TestNewListCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewListOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewListCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.NamespaceFlags, cmdcore.NamespaceFlags{"test-namespace"})
}

func TestNewListCmd_OkMinimum(t *testing.T) {
	realCmd := NewListOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewListCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectReachesExecution()
}
