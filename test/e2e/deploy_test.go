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

func TestBasicDeploy(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName     = "test-basic-deploy-service-name"
		expectedContent = "TestBasicDeploy_Content"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Deploy service", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent,
		})
	})

	logger.Section("Checking if service was added", func() {
		out := knctl.Run([]string{"service", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundService bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == serviceName {
				foundService = true
			}
		}

		if !foundService {
			t.Fatalf("Expected to see sample service in the list of services, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if service details can be seen", func() {
		out := knctl.Run([]string{"service", "show", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if resp.Tables[0].Rows[0]["name"] != serviceName {
			t.Fatalf("Expected to see sample service name in its details, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if service is reachable and presents content", func() {
		curl.WaitForContent(serviceName, expectedContent)
	})

	logger.Section("Check logs of service", func() {
		expectedLogLines := []string{
			"Hello world sample started.",
			"Hello world received a request.",
		}

		out := knctl.Run([]string{"logs", "-s", serviceName})

		for _, line := range expectedLogLines {
			if !strings.Contains(out, line) {
				t.Fatalf("Expected to find log line '%s' in service logs: '%s'", line, out)
			}
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
