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

func TestRoutes(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}

	const (
		routeName    = "test-routes-name"
		serviceName1 = "test-routes-service1-name"
		serviceName2 = "test-routes-service2-name"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"route", "delete", "-n", "default", "--route", routeName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous route with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithOpts([]string{"service", "delete", "-n", "default", "-s", serviceName1}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"service", "delete", "-n", "default", "-s", serviceName2}, RunOpts{AllowError: true})
	})

	defer func() {
		knctl.RunWithOpts([]string{"service", "delete", "-n", "default", "-s", serviceName1}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"service", "delete", "-n", "default", "-s", serviceName2}, RunOpts{AllowError: true})
	}()

	logger.Section("Deploy services that can be routed to", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName1,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + serviceName1,
		})

		knctl.Run([]string{
			"deploy",
			"-s", serviceName2,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + serviceName2,
		})
	})

	logger.Section("Create route", func() {
		knctl.Run([]string{
			"route",
			"create",
			"--route", routeName,
			"-p", serviceName1 + ":latest=50%",
			"-p", serviceName2 + ":latest=50%",
		})
	})

	logger.Section("Checking if route was added", func() {
		out := knctl.Run([]string{"route", "list", "-n", "default", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundRoute bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == routeName {
				foundRoute = true

				traffic := row["traffic"]

				if !strings.Contains(traffic, "50% -> :"+serviceName1) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName1, traffic)
				}
				if !strings.Contains(traffic, "50% -> :"+serviceName2) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName2, traffic)
				}
			}
		}

		if !foundRoute {
			t.Fatalf("Expected to see route in the list of routes, but did not: '%s'", out)
		}
	})

	logger.Section("Check if route directs traffic to both services", func() {
		// TODO figure out why route does not become ready
	})

	logger.Section("Reconfigure route", func() {
		knctl.Run([]string{
			"route",
			"create",
			"--route", routeName,
			"-p", serviceName1 + ":latest=20%",
			"-p", serviceName2 + ":latest=80%",
		})
	})

	logger.Section("Checking if route was reconfigured", func() {
		out := knctl.Run([]string{"route", "list", "-n", "default", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundRoute bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == routeName {
				foundRoute = true

				traffic := row["traffic"]

				if !strings.Contains(traffic, "20% -> :"+serviceName1) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName1, traffic)
				}
				if !strings.Contains(traffic, "80% -> :"+serviceName2) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName2, traffic)
				}
			}
		}

		if !foundRoute {
			t.Fatalf("Expected to see route in the list of routes, but did not: '%s'", out)
		}
	})

	logger.Section("Check if route directs traffic to both services after being reconfigured", func() {
		// TODO figure out why route does not become ready
	})

	logger.Section("Deleting route", func() {
		knctl.Run([]string{"route", "delete", "-n", "default", "--route", routeName})

		out := knctl.Run([]string{"route", "list", "-n", "default", "--json"})
		if strings.Contains(out, routeName) {
			t.Fatalf("Expected to not see route in the list of routes, but was: %s", out)
		}
	})
}
