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

func TestDeployWithBuild(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	knctl := Knctl{t, logger}
	curl := Curl{t, knctl}

	const (
		serviceName      = "test-deploy-with-build-service-name"
		expectedContent1 = "TestDeployWithBuild_ContentV1"
		expectedContent2 = "TestDeployWithBuild_ContentV2"
	)

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	})

	defer func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	}()

	logger.Section("Deploy service v1", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"--git-url", env.BuildGitURL,
			"--git-revision", env.BuildGitRevisionV1,
			"-i", env.BuildImage,
			"--service-account-name", env.BuildServiceAccount,
			"-e", "SIMPLE_MSG=" + expectedContent1,
		})
	})

	logger.Section("Checking if service is reachable and presents content", func() {
		curl.WaitForContent(serviceName, expectedContent1)
	})

	logger.Section("Deploy service v2 with a Git change (new env variable)", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"--git-url", env.BuildGitURL,
			"--git-revision", env.BuildGitRevisionV2,
			"-i", env.BuildImage,
			"--service-account-name", env.BuildServiceAccount,
			"-e", "SIMPLE_MSG_V2=" + expectedContent2,
		})
	})

	logger.Section("Checking if service is reachable and presents content", func() {
		curl.WaitForContent(serviceName, expectedContent2)
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"delete", "service", "-n", "default", "-s", serviceName})

		out := knctl.Run([]string{"list", "services", "-n", "default", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
