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

func TestNewInstallCmd_Ok(t *testing.T) {
	realCmd := NewInstallOptions(nil, NewDepsFactoryImpl(), &KubeconfigFlags{})
	cmd := NewTestCmd(t, NewInstallCmd(realCmd))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-p",
		"-m",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.NodePorts, true)
	DeepEqual(t, realCmd.ExcludeMonitoring, true)
}

func TestNewInstallCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewInstallOptions(nil, NewDepsFactoryImpl(), &KubeconfigFlags{})
	cmd := NewTestCmd(t, NewInstallCmd(realCmd))
	cmd.Execute([]string{
		"--node-ports",
		"--exclude-monitoring",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.NodePorts, true)
	DeepEqual(t, realCmd.ExcludeMonitoring, true)
}

func TestNewInstallCmd_OkMinimum(t *testing.T) {
	realCmd := NewInstallOptions(nil, NewDepsFactoryImpl(), &KubeconfigFlags{})
	cmd := NewTestCmd(t, NewInstallCmd(realCmd))
	cmd.Execute([]string{})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.NodePorts, false)
	DeepEqual(t, realCmd.ExcludeMonitoring, false)
}
