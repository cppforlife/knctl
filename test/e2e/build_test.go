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

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestBuildSuccess(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	const (
		buildName               = "test-build-success-service-name"
		buildDockerSecretName   = buildName + "-docker-secret"
		buildServiceAccountName = buildName + "-service-account"
		expectedBuildOutput     = "Taking snapshot of full filesystem" // coming from kaniko
	)

	logger.Section("Delete previous build with the same name if exists", func() {
		knctl.RunWithOpts([]string{"delete", "build", "-b", buildName}, RunOpts{AllowError: true})
	})

	defer func() {
		knctl.RunWithOpts([]string{"delete", "build", "-b", buildName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", buildDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", buildServiceAccountName}, RunOpts{AllowError: true})
	}()

	logger.Section("Add service account with Docker push secret", func() {
		knctl.RunWithOpts([]string{
			"create",
			"basic-auth-secret",
			"-s", buildDockerSecretName,
			"--docker-hub",
			"-u", env.BuildDockerUsername,
			"-p", env.BuildDockerPassword,
		}, RunOpts{Redact: true})

		knctl.Run([]string{"create", "service-account", "-a", buildServiceAccountName, "-s", buildDockerSecretName})
	})

	logger.Section("Run build and see log output", func() {
		out := knctl.Run([]string{
			"build",
			"-b", buildName,
			"--git-url", env.BuildGitURL,
			"--git-revision", env.BuildGitRevision,
			"-i", env.BuildImage,
			"--service-account-name", buildServiceAccountName,
		})

		// TODO stronger assertion of generated image?
		if !strings.Contains(out, expectedBuildOutput) {
			t.Fatalf("Expected to see kaniko output, but was: %s", out)
		}

		if !strings.Contains(out, env.BuildImage) {
			t.Fatalf("Expected to see image pushed, but was: %s", out)
		}
	})

	logger.Section("Checking if build was added", func() {
		out := knctl.Run([]string{"list", "builds", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundService bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == buildName {
				foundService = true

				if row["succeeded"] != "true" {
					t.Fatalf("Expected build to be marked successful, but was: %#v", row)
				}
			}
		}

		if !foundService {
			t.Fatalf("Expected to see build in the list of builds, but did not: '%s'", out)
		}
	})

	logger.Section("Deleting build", func() {
		knctl.Run([]string{"delete", "build", "-b", buildName})

		out := knctl.Run([]string{"list", "builds", "--json"})
		if strings.Contains(out, buildName) {
			t.Fatalf("Expected to not see build in the list of builds, but was: %s", out)
		}
	})
}

func TestBuildFailed(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	const (
		buildName               = "test-build-failed-service-name"
		buildDockerSecretName   = buildName + "-docker-secret"
		buildServiceAccountName = buildName + "-service-account"
		expectedErrorOuput      = "Unexpected error running git"
	)

	logger.Section("Delete previous build with the same name if exists", func() {
		knctl.RunWithOpts([]string{"delete", "build", "-b", buildName}, RunOpts{AllowError: true})
	})

	defer func() {
		knctl.RunWithOpts([]string{"delete", "build", "-b", buildName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "secret", buildDockerSecretName}, RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", buildServiceAccountName}, RunOpts{AllowError: true})
	}()

	logger.Section("Add service account with Docker push secret", func() {
		knctl.RunWithOpts([]string{
			"create",
			"basic-auth-secret",
			"-s", buildDockerSecretName,
			"--docker-hub",
			"-u", env.BuildDockerUsername,
			"-p", env.BuildDockerPassword,
		}, RunOpts{Redact: true})

		knctl.Run([]string{"create", "service-account", "-a", buildServiceAccountName, "-s", buildDockerSecretName})
	})

	logger.Section("Run build and see it fail", func() {
		out, err := knctl.RunWithOpts([]string{
			"build",
			"-b", buildName,
			"--git-url", "invalid-git-url",
			"--git-revision", "invalid-git-revision",
			"-i", env.BuildImage,
			"--service-account-name", buildServiceAccountName,
		}, RunOpts{AllowError: true})

		if err == nil {
			t.Fatalf("Expected for the command to error")
		}

		// TODO sometimes tailing doesnt pick up output
		// even though if you do kubectl logs -f it shows up
		if !strings.Contains(out, expectedErrorOuput) {
			t.Fatalf("Expected to see error in the log, but was: %s", out)
		}
	})

	defer func() {
		knctl.RunWithOpts([]string{"delete", "build", "-b", buildName}, RunOpts{AllowError: true})
	}()

	logger.Section("Checking if build was added", func() {
		out := knctl.Run([]string{"list", "builds", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundService bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == buildName {
				foundService = true

				if row["succeeded"] != "false" {
					t.Fatalf("Expected build to be marked successful, but was: %#v", row)
				}
			}
		}

		if !foundService {
			t.Fatalf("Expected to see build in the list of builds, but did not: '%s'", out)
		}
	})

	logger.Section("Deleting build", func() {
		knctl.Run([]string{"delete", "build", "-b", buildName})

		out := knctl.Run([]string{"list", "builds", "--json"})
		if strings.Contains(out, buildName) {
			t.Fatalf("Expected to not see build in the list of builds, but was: %s", out)
		}
	})
}
