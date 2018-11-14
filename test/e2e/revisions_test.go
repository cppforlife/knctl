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
	"time"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestRevisions(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName         = "test-revisions-service-name"
		serviceName2        = "test-revisions-service-name2"
		expectedContentRev1 = "TestRevisions_ContentRev1"
		expectedContentRev2 = "TestRevisions_ContentRev2"
		expectedContentRev3 = "TestRevisions_ContentRev3"
	)

	logger.Section("Sleeping...", func() {
		// TODO otherwise 'no upstream healty' error happens
		// somehow caused by previous deploy in other tests
		time.Sleep(20 * time.Second)
	})

	cleanUp := func() {
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName}, RunOpts{AllowError: true})
		knctl.RunWithOpts([]string{"service", "delete", "-s", serviceName2}, RunOpts{AllowError: true})
	}

	logger.Section("Delete previous service with the same name if exists", cleanUp)
	defer cleanUp()

	logger.Section("Deploy revision 1", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev1,
		})
	})

	logger.Section("Checking if revision was added", func() {
		out := knctl.Run([]string{"revision", "list", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 1 {
			t.Fatalf("Expected to see one revision in the list of revisions, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if revision details can be seen", func() {
		out := knctl.Run([]string{"revision", "show", "-r", serviceName + ":latest", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if !strings.Contains(resp.Tables[0].Rows[0]["name"], serviceName) {
			t.Fatalf("Expected to see sample revision name in its details, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if service is reachable and presents content from revision 1", func() {
		curl.WaitForContent(serviceName, expectedContentRev1)
	})

	logger.Section("Deploy revision 2", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev2,
		})
	})

	logger.Section("Checking if revision was added", func() {
		out := knctl.Run([]string{"revision", "list", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 2 {
			t.Fatalf("Expected to see 2 revisions in the list of revisions, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if service is reachable and presents content from revision 2", func() {
		curl.WaitForContent(serviceName, expectedContentRev2)
	})

	logger.Section("Deploy revision 3", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev3,
		})
	})

	logger.Section("Checking if revision was added", func() {
		out := knctl.Run([]string{"revision", "list", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 3 {
			t.Fatalf("Expected to see 3 revisions in the list of revisions, but did not: '%s'", out)
		}
	})

	logger.Section("Checking if service is reachable and presents content from revision 3", func() {
		curl.WaitForContent(serviceName, expectedContentRev3)
	})

	logger.Section("Deleting revision", func() {
		knctl.Run([]string{"revision", "delete", "-r", serviceName + "-00002"}) // TODO better way to find out?
	})

	logger.Section("Checking if revison was deleted", func() {
		out := knctl.Run([]string{"revision", "list", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 2 {
			t.Fatalf("Expected to see 2 revisions in the list of revisions, but did not: '%s'", out)
		}
	})

	logger.Section("Deploy another service", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName2,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev1,
		})
	})

	logger.Section("Checking if revisions from both services can be viewed", func() {
		out := knctl.Run([]string{"revision", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 3 {
			t.Fatalf("Expected to see revisions from both services in the list, but did not: '%s'", out)
		}

		// And if revisions are filtered by service...
		out = knctl.Run([]string{"revision", "list", "-s", serviceName2, "--json"})
		resp = uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 1 {
			t.Fatalf("Expected to see revision from single service in the list, but did not: '%s'", out)
		}
	})
}
