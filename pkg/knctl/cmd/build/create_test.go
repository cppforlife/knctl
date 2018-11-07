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
	"time"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
)

func TestNewCreateCmd_Ok(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-b", "test-build",
		"--git-url", "test-git-url",
		"--git-revision", "test-git-revision",
		"--service-account", "test-service-account",
		"-i", "test-image",
		"--timeout", "1s",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.BuildFlags,
		BuildFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-build"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		cmdcore.GenerateNameFlags{},
		CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:             "test-git-url",
				GitRevision:        "test-git-revision",
				ServiceAccountName: "test-service-account",
				Image:              "test-image",
				Timeout:            1 * time.Second,
			},
		},
	})
}

func TestNewCreateCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--build", "test-build",
		"--git-url", "test-git-url",
		"--git-revision", "test-git-revision",
		"--service-account", "test-service-account",
		"--image", "test-image",
		"--timeout", "1s",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.BuildFlags,
		BuildFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-build"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		cmdcore.GenerateNameFlags{},
		CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:             "test-git-url",
				GitRevision:        "test-git-revision",
				ServiceAccountName: "test-service-account",
				Image:              "test-image",
				Timeout:            1 * time.Second,
			},
		},
	})
}

func TestNewCreateCmd_OkMinimum(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--build", "test-build",
		"--git-url", "test-git-url",
		"--git-revision", "test-git-revision",
		"--image", "test-image",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.BuildFlags,
		BuildFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-build"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		cmdcore.GenerateNameFlags{},
		CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:      "test-git-url",
				GitRevision: "test-git-revision",
				Image:       "test-image",
			},
		},
	})
}

func TestNewCreateCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewConfigFactoryImpl(), cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"build", "image"})
}
