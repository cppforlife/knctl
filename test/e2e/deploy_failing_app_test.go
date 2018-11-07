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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeployFailingApp(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	const (
		serviceName              = "test-d-f-a-service-name"
		pushPullDockerSecretName = serviceName + "-docker-secret"
		pullDockerSecretName     = serviceName + "-p-docker-secret"
		buildServiceAccountName  = serviceName + "-service-account"
		expectedContent1         = "TestDeployBuild_ContentV1"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", pushPullDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", pullDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", buildServiceAccountName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Add service account with Docker push secret", func() {
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
			"-s", pushPullDockerSecretName,
			"-s", pullDockerSecretName,
		})
	})

	cwdPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Expected not to fail getting current directory: %s", err)
	}

	logger.Section("Deploy service that starts but fails right after", func() {
		out, _ := knctl.RunWithOpts([]string{
			"deploy",
			"--tty",
			"-s", serviceName,
			"-d", filepath.Join(cwdPath, "assets", "simple-app-failing-1"),
			"-i", env.BuildPrivateImage,
			"--service-account", buildServiceAccountName,
			"-e", "SIMPLE_MSG=" + expectedContent1,
			"--watch-revision-ready-timeout", "30s",
		}, RunOpts{AllowError: true})

		if !strings.Contains(out, "app-is-exiting") {
			t.Fatalf("Expected to see app failure in the logs: %s", out)
		}
	})

	logger.Section("Deploy service that fails to start due to wrong entrypoint", func() {
		out, _ := knctl.RunWithOpts([]string{
			"deploy",
			"--tty",
			"-s", serviceName,
			"-d", filepath.Join(cwdPath, "assets", "simple-app-failing-2"),
			"-i", env.BuildPrivateImage,
			"--service-account", buildServiceAccountName,
			"-e", "SIMPLE_MSG=" + expectedContent1,
			"--watch-revision-ready-timeout", "30s",
		}, RunOpts{AllowError: true})

		if !strings.Contains(out, "stat /wrong-app: no such file or directory") {
			t.Fatalf("Expected to see error message about app: %s", out)
		}

		if !strings.Contains(out, fmt.Sprintf("Revision '%s-00002' did not became ready", serviceName)) {
			t.Fatalf("Expected to see revision did not become ready: %s", out)
		}
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"service", "delete", "-s", serviceName})

		out := knctl.Run([]string{"service", "list", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
