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
	curl := Curl{t, knctl}

	const (
		routeName    = "test-routes-name"
		serviceName1 = "test-routes-service1-name"
		serviceName2 = "test-routes-service2-name"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"route", "delete", "--route", routeName}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName1}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName2}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous route and services if exists", cleanUp)
	defer cleanUp()

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
			"rollout",
			"--route", routeName,
			"-p", serviceName1 + ":latest=50%",
			"-p", serviceName2 + ":latest=50%",
		})
	})

	logger.Section("Checking if route was added", func() {
		out := knctl.Run([]string{"route", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundRoute bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == routeName {
				foundRoute = true

				traffic := row["traffic"]

				if !strings.Contains(traffic, "50% -> "+serviceName1) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName1, traffic)
				}
				if !strings.Contains(traffic, "50% -> "+serviceName2) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName2, traffic)
				}
			}
		}

		if !foundRoute {
			t.Fatalf("Expected to see route in the list of routes, but did not: '%s'", out)
		}
	})

	logger.Section("Check if route directs traffic to both services", func() {
		// Make sure route returns both first
		curl.WaitForRouteContent(routeName, serviceName1)
		curl.WaitForRouteContent(routeName, serviceName2)

		counts := curl.RouteContentCounts(routeName, 100, []string{serviceName1, serviceName2})
		if len(counts) != 2 {
			t.Fatalf("Expected route to point to only two services and received two types of responses, counts: %#v", counts)
		}

		for _, name := range []string{serviceName1, serviceName2} {
			v := counts[name]
			if v < 40 && v > 60 {
				t.Fatalf("Expected route to split traffic equally between two services, counts: %#v", counts)
			}
		}
	})

	logger.Section("Reconfigure route", func() {
		knctl.Run([]string{
			"rollout",
			"--route", routeName,
			"-p", serviceName1 + ":latest=20%",
			"-p", serviceName2 + ":latest=80%",
		})
	})

	logger.Section("Checking if route was reconfigured", func() {
		out := knctl.Run([]string{"route", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundRoute bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == routeName {
				foundRoute = true

				traffic := row["traffic"]

				if !strings.Contains(traffic, "20% -> "+serviceName1) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName1, traffic)
				}
				if !strings.Contains(traffic, "80% -> "+serviceName2) {
					t.Fatalf("Expected route to point to '%s', but was '%s'", serviceName2, traffic)
				}
			}
		}

		if !foundRoute {
			t.Fatalf("Expected to see route in the list of routes, but did not: '%s'", out)
		}
	})

	logger.Section("Check if route directs traffic to both services after being reconfigured", func() {
		// Make sure route returns both first
		curl.WaitForRouteContent(routeName, serviceName1)
		curl.WaitForRouteContent(routeName, serviceName2)

		counts := curl.RouteContentCounts(routeName, 100, []string{serviceName1, serviceName2})
		if len(counts) != 2 {
			t.Fatalf("Expected route to point to only two services and received two types of responses, counts: %#v", counts)
		}

		v1 := counts[serviceName1]
		if v1 < 10 && v1 > 30 {
			t.Fatalf("Expected route to split 20%% of traffic to service1, counts: %#v", counts)
		}

		v2 := counts[serviceName2]
		if v2 < 70 && v2 > 90 {
			t.Fatalf("Expected route to split 80%% of traffic to service2, counts: %#v", counts)
		}
	})
}
