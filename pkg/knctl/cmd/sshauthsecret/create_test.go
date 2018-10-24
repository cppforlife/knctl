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

package sshauthsecret_test

import (
	"testing"

	. "github.com/cppforlife/knctl/pkg/knctl/cmd"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/sshauthsecret"
)

func TestNewCreateCmd_Ok(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-s", "test-secret",
		"--url", "test-url",
		"--private-key", "test-private-key-pem",
		"--private-key-path", "test-private-key-path",
		"--known-hosts", "test-known-hosts",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		URL:            "test-url",
		PrivateKey:     "test-private-key-pem",
		PrivateKeyPath: "test-private-key-path",
		KnownHosts:     "test-known-hosts",
	})
}

func TestNewCreateCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--url", "test-url",
		"--private-key", "test-private-key-pem",
		"--known-hosts", "test-known-hosts",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		URL:        "test-url",
		PrivateKey: "test-private-key-pem",
		KnownHosts: "test-known-hosts",
	})
}

func TestNewCreateCmd_OkGithub(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--private-key", "test-private-key-pem",
		"--known-hosts", "test-known-hosts",
		"--github",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		cmdflags.SecretFlags{cmdcore.NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:       "",
		URL:        "",
		PrivateKey: "test-private-key-pem",
		KnownHosts: "test-known-hosts",
		Github:     true,
	})

	err := realCmd.CreateFlags.BackfillTypeAndURL()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	DeepEqual(t, realCmd.CreateFlags, CreateFlags{
		Type:       "git",
		URL:        "github.com",
		PrivateKey: "test-private-key-pem",
		KnownHosts: "test-known-hosts",
		Github:     true,
	})
}

func TestNewCreateCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateOptions(nil, cmdcore.NewDepsFactory())
	cmd := NewTestCmd(t, NewCreateCmd(realCmd, cmdcore.FlagsFactory{}))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"secret"})
}
