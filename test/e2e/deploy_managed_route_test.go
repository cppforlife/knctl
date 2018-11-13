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

func TestDeployManagedRouteLaterDeploy(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName      = "test-deploy-managed-route-service-name"
		routeName        = serviceName // same as service
		expectedContent1 = "TestDeployCustomRoute_Content1"
		expectedContent2 = "TestDeployCustomRoute_Content2"
		expectedContent3 = "TestDeployCustomRoute_Content3"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"route", "delete", "--route", routeName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Deploy service and manually create route", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent1,
			"--managed-route=false",
		})

		// After this deploy no associated route exists

		knctl.Run([]string{
			"rollout",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

		curl.WaitForRouteContent(routeName, expectedContent1)
	})

	logger.Section("Deploy service and manually route to latest revision", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent2,
			"--managed-route=false",
		})

		// still should be returning content 1
		curl.WaitForRouteContent(routeName, expectedContent1)

		knctl.Run([]string{
			"rollout",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

		// switches to latest content
		curl.WaitForRouteContent(routeName, expectedContent2)
	})

	logger.Section("Deploy service and revert to managed route", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent3,
		})

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

func TestDeployManagedRouteFirstDeploy(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName      = "test-deploy-managed-route-service-name"
		routeName        = serviceName // same as service
		expectedContent1 = "TestDeployCustomRoute_Content1"
		expectedContent2 = "TestDeployCustomRoute_Content2"
		expectedContent3 = "TestDeployCustomRoute_Content3"
	)

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"route", "delete", "--route", routeName}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Deploy service with a managed route at first", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent1,
		})

		curl.WaitForRouteContent(routeName, expectedContent1)
	})

	logger.Section("Deploy service with unmanaged route", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent2,
			"--managed-route=false",
		})

		// will be routing to content 2 because previously route was managed
		// TODO somewhat awkward user experience as until rollout command is called
		// existing route will keep on tracking latest configuration
		curl.WaitForRouteContent(routeName, expectedContent2)

		knctl.Run([]string{
			"rollout",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})
	})

	logger.Section("Deploy service and rollout manually", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContent3,
			"--managed-route=false",
		})

		curl.WaitForRouteContent(serviceName, expectedContent2)

		knctl.Run([]string{
			"rollout",
			"--route", routeName,
			"-p", serviceName + ":latest=100%",
		})

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
