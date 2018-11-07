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

package service_test

import (
	"testing"
	"time"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	cmdbld "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/service"
)

func TestNewDeployCmd_Ok(t *testing.T) {
	realCmd := NewDeployOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewDeployCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-s", "test-service",
		"--git-url", "test-git-url",
		"--git-revision", "test-git-revision",
		"--service-account", "test-service-account",
		"-i", "test-image",
		"-e", "key1=val1",
		"-e", "key2=val2",
		"--build-timeout", "1s",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		cmdflags.ServiceFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-service"})

	DeepEqual(t, realCmd.DeployFlags, DeployFlags{
		BuildCreateArgsFlags: cmdbld.CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:             "test-git-url",
				GitRevision:        "test-git-revision",
				ServiceAccountName: "test-service-account",
				Timeout:            1 * time.Second,
			},
		},
		Image:              "test-image",
		Env:                []string{"key1=val1", "key2=val2"},
		WatchRevisionReady: true,
		WatchPodLogs:       true,
	})
}

func TestNewDeployCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewDeployOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewDeployCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--service", "test-service",
		"--git-url", "test-git-url",
		"--git-revision", "test-git-revision",
		"--service-account", "test-service-account",
		"--image", "test-image",
		"--env", "key1=val1",
		"--env", "key2=val2",
		"--build-timeout", "1s",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		cmdflags.ServiceFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-service"})

	DeepEqual(t, realCmd.DeployFlags, DeployFlags{
		BuildCreateArgsFlags: cmdbld.CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:             "test-git-url",
				GitRevision:        "test-git-revision",
				ServiceAccountName: "test-service-account",
				Timeout:            1 * time.Second,
			},
		},
		Image:              "test-image",
		Env:                []string{"key1=val1", "key2=val2"},
		WatchRevisionReady: true,
		WatchPodLogs:       true,
	})
}

func TestNewDeployCmd_OkMinimum(t *testing.T) {
	realCmd := NewDeployOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewDeployCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--service", "test-service",
		"--image", "test-image",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		cmdflags.ServiceFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-service"})

	DeepEqual(t, realCmd.DeployFlags, DeployFlags{
		Image:              "test-image",
		WatchRevisionReady: true,
		WatchPodLogs:       true,
	})
}

func TestNewDeployCmd_WatchingDisabled(t *testing.T) {
	realCmd := NewDeployOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewDeployCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--service", "test-service",
		"--image", "test-image",
		"--watch-revision-ready=false",
		"--watch-pod-logs=false",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.ServiceFlags,
		cmdflags.ServiceFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-service"})

	DeepEqual(t, realCmd.DeployFlags, DeployFlags{
		Image:              "test-image",
		WatchRevisionReady: false,
		WatchPodLogs:       false,
	})
}

func TestNewDeployCmd_RequiredFlags(t *testing.T) {
	realCmd := NewDeployOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewDeployCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"image", "service"})
}
