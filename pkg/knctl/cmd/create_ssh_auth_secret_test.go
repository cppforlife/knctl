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

func TestNewCreateSSHAuthSecretCmd_Ok(t *testing.T) {
	realCmd := NewCreateSSHAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateSSHAuthSecretCmd(realCmd))
	cmd.ExpectBasicConfig()
	cmd.Execute([]string{
		"-n", "test-namespace",
		"-s", "test-secret",
		"--url", "test-url",
		"--private-key", "test-private-key-pem",
		"--known-hosts", "test-known-hosts",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.SSHAuthSecretCreateFlags, SSHAuthSecretCreateFlags{
		URL:        "test-url",
		PrivateKey: "test-private-key-pem",
		KnownHosts: "test-known-hosts",
	})
}

func TestNewCreateSSHAuthSecretCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateSSHAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateSSHAuthSecretCmd(realCmd))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--url", "test-url",
		"--private-key", "test-private-key-pem",
		"--known-hosts", "test-known-hosts",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.SSHAuthSecretCreateFlags, SSHAuthSecretCreateFlags{
		URL:        "test-url",
		PrivateKey: "test-private-key-pem",
		KnownHosts: "test-known-hosts",
	})
}

func TestNewCreateSSHAuthSecretCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateSSHAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateSSHAuthSecretCmd(realCmd))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"namespace", "private-key", "secret", "url"})
}
