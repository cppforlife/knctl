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
	"reflect"
	"sort"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestDefaultRevisionTags(t *testing.T) {
	logger := Logger{}
	knctl := Knctl{t, logger}

	const (
		serviceName         = "test-default-revisions-tags-service-name"
		expectedContentRev1 = "TestRevisions_ContentRev1"
		expectedContentRev2 = "TestRevisions_ContentRev2"
		expectedContentRev3 = "TestRevisions_ContentRev3"
	)

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	})

	defer func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	}()

	logger.Section("Deploy first revision", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev1,
		})
	})

	expectNumberOfRevisions := func(expectedRows int) []map[string]string {
		out := knctl.Run([]string{"list", "revisions", "-n", "default", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != expectedRows {
			t.Fatalf("Expected to see one revision in the list of revisions, but did not: '%s'", out)
		}

		return resp.Tables[0].Rows
	}

	logger.Section("Checking that tags 'latest' and 'previous' are on the first revision", func() {
		latestRow := expectNumberOfRevisions(1)[0]
		tags := strings.Split(latestRow["tags"], "\n")
		sort.Strings(tags)

		if len(tags) != 2 {
			t.Fatalf("Expected 2 tags on the revision, but was '%#v'", tags)
		}

		if !reflect.DeepEqual(tags, []string{"latest", "previous"}) {
			t.Fatalf("Expected latest and previous tags on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Deploy second revision", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev2,
		})
	})

	logger.Section("Checking that tags 'latest' is on the second revision", func() {
		latestRow := expectNumberOfRevisions(2)[0]
		tags := strings.Split(latestRow["tags"], "\n")

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag on the revision, but was '%#v'", tags)
		}

		if !reflect.DeepEqual(tags, []string{"latest"}) {
			t.Fatalf("Expected tag 'latest' on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Checking that tags 'previous' is on the first revision", func() {
		oldestRow := expectNumberOfRevisions(2)[1]
		tags := strings.Split(oldestRow["tags"], "\n")

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag on the revision, but was '%#v'", tags)
		}

		if !reflect.DeepEqual(tags, []string{"previous"}) {
			t.Fatalf("Expected tag 'previous' on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Deploy third revision", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev3,
		})
	})

	logger.Section("Checking that tags 'latest' is on the third revision", func() {
		latestRow := expectNumberOfRevisions(3)[0]
		tags := strings.Split(latestRow["tags"], "\n")

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag on the revision, but was '%#v'", tags)
		}

		if !reflect.DeepEqual(tags, []string{"latest"}) {
			t.Fatalf("Expected tag 'latest' on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Checking that tags 'previous' is on the second revision", func() {
		prevRow := expectNumberOfRevisions(3)[1]
		tags := strings.Split(prevRow["tags"], "\n")

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag on the revision, but was '%#v'", tags)
		}

		if !reflect.DeepEqual(tags, []string{"previous"}) {
			t.Fatalf("Expected tag 'previous' on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Checking that no tags is on the first revision", func() {
		oldestRow := expectNumberOfRevisions(3)[2]
		tags := oldestRow["tags"]

		if len(tags) != 0 {
			t.Fatalf("Expected 0 tags on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"delete", "service", "-n", "default", "-s", serviceName})

		out := knctl.Run([]string{"list", "services", "-n", "default", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}

func TestTagRevisions(t *testing.T) {
	logger := Logger{}
	knctl := Knctl{t, logger}

	const (
		serviceName         = "test-tag-revisions-service-name"
		expectedContentRev1 = "TestRevisions_ContentRev1"
		expectedContentRev2 = "TestRevisions_ContentRev2"
		expectedContentRev3 = "TestRevisions_ContentRev3"
	)

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	})

	defer func() {
		knctl.RunWithErr([]string{"delete", "service", "-n", "default", "-s", serviceName})
	}()

	logger.Section("Deploy two revisions", func() {
		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev1,
		})

		knctl.Run([]string{
			"deploy",
			"-n", "default",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev2,
		})
	})

	expectNumberOfRevisions := func(expectedRows int) []map[string]string {
		out := knctl.Run([]string{"list", "revisions", "-n", "default", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != expectedRows {
			t.Fatalf("Expected to see one revision in the list of revisions, but did not: '%s'", out)
		}

		return resp.Tables[0].Rows
	}

	extractTags := func(row map[string]string) []string {
		tags := strings.Split(row["tags"], "\n")
		sort.Strings(tags)
		return tags
	}

	logger.Section("Check that revision can be tagged", func() {
		firstRow := expectNumberOfRevisions(2)[0]
		knctl.Run([]string{"tag", "revision", "-t", "tag1", "-t", "tag2", "-n", "default", "-r", firstRow["name"], "--json"})

		tags := extractTags(expectNumberOfRevisions(2)[0])
		if !reflect.DeepEqual(tags, []string{"latest", "tag1", "tag2"}) {
			t.Fatalf("Expected latest and previous tags on the revision, but was '%#v'", tags)
		}
	})

	logger.Section("Check that revision can be re-tagged", func() {
		lastRow := expectNumberOfRevisions(2)[1]
		knctl.Run([]string{"tag", "revision", "-t", "tag2", "-n", "default", "-r", lastRow["name"], "--json"})

		tags := extractTags(expectNumberOfRevisions(2)[0])
		if !reflect.DeepEqual(tags, []string{"latest", "tag1"}) {
			t.Fatalf("Expected tag 'tag2' was removed from latest revision, but tags were '%#v'", tags)
		}

		tags = extractTags(expectNumberOfRevisions(2)[1])
		if !reflect.DeepEqual(tags, []string{"previous", "tag2"}) {
			t.Fatalf("Expected tag 'tag2' was added to the earlier revision, but was '%#v'", tags)
		}
	})

	logger.Section("Check that revision can be untagged", func() {
		firstRow := expectNumberOfRevisions(2)[0]
		knctl.Run([]string{"untag", "revision", "-t", "latest", "-n", "default", "-r", firstRow["name"], "--json"})

		tags := extractTags(expectNumberOfRevisions(2)[0])
		if !reflect.DeepEqual(tags, []string{"tag1"}) {
			t.Fatalf("Expected tag 'latest' was removed, but tags were '%#v'", tags)
		}
	})

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"delete", "service", "-n", "default", "-s", serviceName})

		out := knctl.Run([]string{"list", "services", "-n", "default", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}
