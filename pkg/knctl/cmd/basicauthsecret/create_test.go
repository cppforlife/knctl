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

package basicauthsecret_test

import (
	"testing"

	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/basicauthsecret"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
)

func TestNewCreateCmd_Ok(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-s", "test-secret",
		"--type", "test-type",
		"--url", "test-url",
		"-u", "test-username",
		"-p", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:     "test-type",
		URL:      "test-url",
		Username: "test-username",
		Password: "test-password",
	})
}

func TestNewCreateCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--type", "test-type",
		"--url", "test-url",
		"--username", "test-username",
		"--password", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:     "test-type",
		URL:      "test-url",
		Username: "test-username",
		Password: "test-password",
	})
}

func TestNewCreateCmd_OkDockerHub(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--docker-hub",
		"--username", "test-username",
		"--password", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:      "",
		URL:       "",
		Username:  "test-username",
		Password:  "test-password",
		DockerHub: true,
	})

	err := realCmd.CreateFlags.BackfillTypeAndURL()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:      "docker",
		URL:       "https://index.docker.io/v1/",
		Username:  "test-username",
		Password:  "test-password",
		DockerHub: true,
	})
}

func TestNewCreateCmd_OkGCR(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--gcr",
		"--username", "test-username",
		"--password", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:     "",
		URL:      "",
		Username: "test-username",
		Password: "test-password",
		GCR:      true,
	})

	err := realCmd.CreateFlags.BackfillTypeAndURL()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:     "docker",
		URL:      "https://gcr.io",
		Username: "test-username",
		Password: "test-password",
		GCR:      true,
	})
}

func TestNewCreateCmd_OkForPulling(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--type", "test-type",
		"--url", "test-url",
		"--username", "test-username",
		"--password", "test-password",
		"--for-pulling",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:       "test-type",
		URL:        "test-url",
		Username:   "test-username",
		Password:   "test-password",
		ForPulling: true,
	})
}

func TestNewCreateCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"password", "secret", "username"})
}
