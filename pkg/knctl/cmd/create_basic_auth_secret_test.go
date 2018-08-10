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

func TestNewCreateBasicAuthSecretCmd_Ok(t *testing.T) {
	realCmd := NewCreateBasicAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateBasicAuthSecretCmd(realCmd))
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
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.BasicAuthSecretCreateFlags, BasicAuthSecretCreateFlags{
		Type:     "test-type",
		URL:      "test-url",
		Username: "test-username",
		Password: "test-password",
	})
}

func TestNewCreateBasicAuthSecretCmd_OkLongFlagNames(t *testing.T) {
	realCmd := NewCreateBasicAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateBasicAuthSecretCmd(realCmd))
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
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.BasicAuthSecretCreateFlags, BasicAuthSecretCreateFlags{
		Type:     "test-type",
		URL:      "test-url",
		Username: "test-username",
		Password: "test-password",
	})
}

func TestNewCreateBasicAuthSecretCmd_OkDockerHub(t *testing.T) {
	realCmd := NewCreateBasicAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateBasicAuthSecretCmd(realCmd))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--docker-hub",
		"--username", "test-username",
		"--password", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.BasicAuthSecretCreateFlags, BasicAuthSecretCreateFlags{
		Type:      "",
		URL:       "",
		Username:  "test-username",
		Password:  "test-password",
		DockerHub: true,
	})
}

func TestNewCreateBasicAuthSecretCmd_OkGCR(t *testing.T) {
	realCmd := NewCreateBasicAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateBasicAuthSecretCmd(realCmd))
	cmd.Execute([]string{
		"--namespace", "test-namespace",
		"--secret", "test-secret",
		"--gcr",
		"--username", "test-username",
		"--password", "test-password",
	})
	cmd.ExpectReachesExecution()

	DeepEqual(t, realCmd.SecretFlags,
		SecretFlags{NamespaceFlags{"test-namespace"}, "test-secret"})

	DeepEqual(t, realCmd.BasicAuthSecretCreateFlags, BasicAuthSecretCreateFlags{
		Type:     "",
		URL:      "",
		Username: "test-username",
		Password: "test-password",
		GCR:      true,
	})
}

func TestNewCreateBasicAuthSecretCmd_RequiredFlags(t *testing.T) {
	realCmd := NewCreateBasicAuthSecretOptions(nil, NewDepsFactoryImpl())
	cmd := NewTestCmd(t, NewCreateBasicAuthSecretCmd(realCmd))
	cmd.Execute([]string{})
	cmd.ExpectRequiredFlags([]string{"namespace", "password", "secret", "username"})
}
