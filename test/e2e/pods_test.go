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

func TestPods(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName      = "test-pods-service-name"
		expectedContent1 = "TestBasicDeploy_Content1"
		expectedContent2 = "TestBasicDeploy_Content2"
	)

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithOpts([]string{"delete", "service", "-s", serviceName}, RunOpts{AllowError: true})
	})

	defer func() {
		knctl.RunWithOpts([]string{"delete", "service", "-s", serviceName}, RunOpts{AllowError: true})
	}()

	logger.Section("Deploy service", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent1,
		})

		curl.WaitForContent(serviceName, expectedContent1)
	})

	logger.Section("Deploy additional revision", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent2,
		})

		curl.WaitForContent(serviceName, expectedContent2)
	})

	logger.Section("Check listing of pods", func() {
		out := knctl.Run([]string{"list", "pods", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 2 {
			t.Fatalf("Expected to find 2 pods, but was '%s'", out)
		}

		pod0 := resp.Tables[0].Rows[0]
		pod1 := resp.Tables[0].Rows[1]

		if len(pod0["revision"]) == 0 || len(pod1["revision"]) == 0 {
			t.Fatalf("Expected pods to not have empty revisions, but was '%s'", out)
		}
		if len(pod0["name"]) == 0 || len(pod1["name"]) == 0 {
			t.Fatalf("Expected pods to not have empty names, but was '%s'", out)
		}
		if pod0["revision"] == pod1["revision"] {
			t.Fatalf("Expected 2 pods to be from different revisions, but was '%s'", out)
		}
		if pod0["name"] == pod1["name"] {
			t.Fatalf("Expected 2 pods to be have two different names, but was '%s'", out)
		}
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"delete", "service", "-s", serviceName})

		out := knctl.Run([]string{"list", "services", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
