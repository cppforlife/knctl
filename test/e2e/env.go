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

	BuildPrivateGit EnvBuildPrivateGit

	BuildPublicImage    string // push with auth, pull w/o auth
	BuildPrivateImage   string // push and pull requires auth
	BuildDockerUsername string
	BuildDockerPassword string
}

type EnvBuildPrivateGit struct {
	URL        string // push and pull requires auth
	SSHPullKey string // key for pulling
	Revision   string
	RevisionV1 string
	RevisionV2 string
}

func BuildEnv(t *testing.T) Env {
	env := Env{
		Namespace: os.Getenv("KNCTL_E2E_NAMESPACE"),

		BuildGitURL:        os.Getenv("KNCTL_E2E_BUILD_GIT_URL"),
		BuildGitRevision:   os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION"),
		BuildGitRevisionV1: os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION_V1"), // See deploy_with_build_test.go for usage
		BuildGitRevisionV2: os.Getenv("KNCTL_E2E_BUILD_GIT_REVISION_V2"),

		// See deploy_build_private_git_private_image_test.go for usage
		BuildPrivateGit: EnvBuildPrivateGit{
			URL:        os.Getenv("KNCTL_E2E_BUILD_PRIVATE_GIT_URL"),
			SSHPullKey: os.Getenv("KNCTL_E2E_BUILD_PRIVATE_GIT_SSH_PULL_KEY"),
			Revision:   os.Getenv("KNCTL_E2E_BUILD_PRIVATE_GIT_REVISION"),
			RevisionV1: os.Getenv("KNCTL_E2E_BUILD_PRIVATE_GIT_REVISION_V1"),
			RevisionV2: os.Getenv("KNCTL_E2E_BUILD_PRIVATE_GIT_REVISION_V2"),
		},

		BuildPublicImage:    os.Getenv("KNCTL_E2E_BUILD_PUBLIC_IMAGE"),
		BuildPrivateImage:   os.Getenv("KNCTL_E2E_BUILD_PRIVATE_IMAGE"),
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

	if len(e.BuildPrivateGit.URL) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateGit.URL to be non-empty")
	}
	if len(e.BuildPrivateGit.SSHPullKey) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateGit.SSHPullKey to be non-empty")
	}
	if len(e.BuildPrivateGit.Revision) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateGit.Revision to be non-empty")
	}
	if len(e.BuildPrivateGit.RevisionV1) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateGit.RevisionV1 to be non-empty")
	}
	if len(e.BuildPrivateGit.RevisionV2) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateGit.RevisionV2 to be non-empty")
	}

	if len(e.BuildPublicImage) == 0 {
		errStrs = append(errStrs, "Expected BuildPublicImage to be non-empty")
	}
	if len(e.BuildPrivateImage) == 0 {
		errStrs = append(errStrs, "Expected BuildPrivateImage to be non-empty")
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
