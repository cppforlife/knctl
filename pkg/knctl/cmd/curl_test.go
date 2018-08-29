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
	"testing"

	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
)

func TestNewCurlCmd_Ok(t *testing.T) {
	realCmd := NewCurlOptions(nil, newDepsFactory())
	cmd := NewTestCmd(t, NewCurlCmd(realCmd, FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-s", "test-service",
		"-p", "1234",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		ServiceFlags{NamespaceFlags{"test-namespace"}, "test-service"})
	DeepEqual(t, realCmd.CurlPortFlags, CurlPortFlags{Port: int32(1234)})
}

func TestNewCurlCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCurlOptions(nil, newDepsFactory())
	cmd := NewTestCmd(t, NewCurlCmd(realCmd, FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--service", "test-service",
		"--port", "1234",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		ServiceFlags{NamespaceFlags{"test-namespace"}, "test-service"})
	DeepEqual(t, realCmd.CurlPortFlags, CurlPortFlags{Port: int32(1234)})
}

func TestNewCurlCmd_OkMinimum(t *testing.T) {
	realCmd := NewCurlOptions(nil, newDepsFactory())
	cmd := NewTestCmd(t, NewCurlCmd(realCmd, FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--service", "test-service",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		ServiceFlags{NamespaceFlags{"test-namespace"}, "test-service"})
	DeepEqual(t, realCmd.CurlPortFlags, CurlPortFlags{Port: int32(80)})
}

func TestNewCurlCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCurlOptions(nil, newDepsFactory())
	cmd := NewTestCmd(t, NewCurlCmd(realCmd, FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"service"})
}
