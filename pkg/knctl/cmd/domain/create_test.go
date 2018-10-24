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

package domain_test

import (
	"testing"

	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/domain"
)

func TestNewCreateCmd_Ok(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-d", "test-domain",
		"--default",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.Domain, "test-domain")
	DeepEqual(t, realCmd.Default, true)
}

func TestNewCreateCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--domain", "test-domain",
		"--default",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.Domain, "test-domain")
	DeepEqual(t, realCmd.Default, true)
}

func TestNewCreateCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"default", "domain"})
}
