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
	"os"
	"strings"
	"testing"
)

type Env struct {
	Namespace string

	BuildGitURL        string
	BuildGitRevision   string
	BuildGitRevisionV1 string
	BuildGitRevisionV2 string

	BuildImage          string
	BuildDockerUsername string
	BuildDockerPassword string
}

func BuildEnv(t *testing.T) Env {
	env := Env{
		Namespace: os.Getenv("KNCTL_E2E_NAMESPACE"),

		BuildGitURL:        os.Getenv("KNCTL_E2E_BUILD_GIT_URL"),
		BuildGitRevision:   os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION"),
		BuildGitRevisionV1: os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION_V1"), // See deploy_with_build_test.go for usage
		BuildGitRevisionV2: os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION_V2"),

		BuildImage:          os.Getenv("KNCTL_E2E_BUILD_IMAGE"),
		BuildDockerUsername: os.Getenv("KNCTL_E2E_BUILD_DOCKER_USERNAME"),
		BuildDockerPassword: os.Getenv("KNCTL_E2E_BUILD_DOCKER_PASSWORD"),
	}
	env.Validate(t)
	return env
}

func (e Env) Validate(t *testing.T) {
	errStrs := []string{}

	if len(e.Namespace) == 0 {
		errStrs = append(errStrs, "Expected Namespace to be non-empty")
	}

	if len(e.BuildGitURL) == 0 {
		errStrs = append(errStrs, "Expected BuildGitURL to be non-empty")
	}
	if len(e.BuildGitRevision) == 0 {
		errStrs = append(errStrs, "Expected BuildGitRevision to be non-empty")
	}
	if len(e.BuildGitRevisionV1) == 0 {
		errStrs = append(errStrs, "Expected BuildGitRevisionV1 to be non-empty")
	}
	if len(e.BuildGitRevisionV2) == 0 {
		errStrs = append(errStrs, "Expected BuildGitRevisionV2 to be non-empty")
	}

	if len(e.BuildImage) == 0 {
		errStrs = append(errStrs, "Expected BuildImage to be non-empty")
	}
	if len(e.BuildDockerUsername) == 0 {
		errStrs = append(errStrs, "Expected BuildDockerUsername to be non-empty")
	}
	if len(e.BuildDockerPassword) == 0 {
		errStrs = append(errStrs, "Expected BuildDockerPassword to be non-empty")
	}

	if len(errStrs) > 0 {
		t.Fatalf("%s", strings.Join(errStrs, "\n"))
	}
}
