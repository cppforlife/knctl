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

func TestDeployCustomRoute(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName      = "test-deploy-custom-route-service-name"
		routeName        = serviceName + "-route"
		expectedContent1 = "TestDeployCustomRoute_Content1"
		expectedContent2 = "TestDeployCustomRoute_Content2"
		expectedContent3 = "TestDeployCustomRoute_Content3"
		expectedContent4 = "TestDeployCustomRoute_Content4"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"route", "delete", "--route", routeName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Deploy service with a custom route", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent1,
			"--custom-route",
		})

		knctl.Run([]string{
			"route", "create",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

		curl.WaitForRouteContent(routeName, expectedContent1)
	})

	logger.Section("Deploy service with a custom route", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent1,
			"--custom-route",
		})

		curl.WaitForRouteContent(routeName, expectedContent1)

		knctl.Run([]string{
			"route", "create",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

		curl.WaitForRouteContent(routeName, expectedContent2)
	})

	logger.Section("Deploy service with 2", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent2,
		})

		curl.WaitForContent(serviceName, expectedContent2)
	})

	logger.Section("Deploy service with 3", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent3,
			"--custom-route",
		})

		knctl.Run([]string{
			"route", "create",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

		// default service route will point to latest version
		curl.WaitForContent(serviceName, expectedContent3)
		// custom route is pointing a new revision
		curl.WaitForRouteContent(serviceName, expectedContent3)
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"service", "delete", "-s", serviceName})

		out := knctl.Run([]string{"service", "list", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
