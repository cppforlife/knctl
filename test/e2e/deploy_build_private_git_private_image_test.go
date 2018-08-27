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

package e2e

import (
	"strings"
	"testing"
)

func TestDeployBuildPrivateGitPrivateImage(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName              = "test-d-b-p-i-p-g-service-name"
		pullGitSecretName        = serviceName + "-p-git-secret"
		pushPullDockerSecretName = serviceName + "-docker-secret"
		pullDockerSecretName     = serviceName + "-p-docker-secret"
		buildServiceAccountName  = serviceName + "-service-account"
		expectedContent1         = "TestDeployBuild_ContentV1"
		expectedContent2         = "TestDeployBuild_ContentV2"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", pullGitSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", pushPullDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", pullDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", buildServiceAccountName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Add service account with Docker push secret", func() {
		if !strings.Contains(env.BuildPrivateGit.URL, "github.com") {
			t.Fatalf("Expected private Git URL '%s' to be github.com URL", env.BuildPrivateGit.URL)
		}

		knctl.RunWithOpts([]string{
			"ssh-auth-secret",
			"create",
			"-s", pullGitSecretName,
			"--github",
			"--private-key", env.BuildPrivateGit.SSHPullKey,
		}, RunOpts{Redact: true})

		knctl.RunWithOpts([]string{
			"basic-auth-secret",
			"create",
			"-s", pushPullDockerSecretName,
			"--docker-hub",
			"-u", env.BuildDockerUsername,
			"-p", env.BuildDockerPassword,
		}, RunOpts{Redact: true})

		knctl.RunWithOpts([]string{
			"basic-auth-secret",
			"create",
			"-s", pullDockerSecretName,
			"--docker-hub",
			"-u", env.BuildDockerUsername,
			"-p", env.BuildDockerPassword,
			"--for-pulling",
		}, RunOpts{Redact: true})

		knctl.Run([]string{
			"service-account",
			"create",
			"-a", buildServiceAccountName,
			"-s", pullGitSecretName,
			"-s", pushPullDockerSecretName,
			"-s", pullDockerSecretName,
		})
	})

	logger.Section("Deploy service v1", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"--git-url", env.BuildPrivateGit.URL,
			"--git-revision", env.BuildPrivateGit.RevisionV1,
			"-i", env.BuildPrivateImage,
			"--service-account", buildServiceAccountName,
			"-e", "SIMPLE_MSG=" + expectedContent1,
		})
	})

	logger.Section("Checking if service is reachable and presents content", func() {
		curl.WaitForContent(serviceName, expectedContent1)
	})

	logger.Section("Deploy service v2 with a Git change (new env variable)", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"--git-url", env.BuildPrivateGit.URL,
			"--git-revision", env.BuildPrivateGit.RevisionV2,
			"-i", env.BuildPrivateImage,
			"--service-account", buildServiceAccountName,
			"-e", "SIMPLE_MSG_V2=" + expectedContent2,
		})
	})

	logger.Section("Checking if service is reachable and presents content", func() {
		curl.WaitForContent(serviceName, expectedContent2)
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"service", "delete", "-s", serviceName})

		out := knctl.Run([]string{"service", "list", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
