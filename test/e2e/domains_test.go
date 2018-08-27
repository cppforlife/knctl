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
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestDomains(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}

	newDomainsSelection := map[string]struct{}{
		"my-domain.test":       struct{}{},
		"my-other-domain.test": struct{}{},
	}

	logger.Section("Checking if at least one domain is available by default", func() {
		out, _ := knctl.RunWithOpts([]string{"domain", "list", "--json"}, RunOpts{NoNamespace: true})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		var foundDefaultDomain bool

		for _, row := range resp.Tables[0].Rows {
			delete(newDomainsSelection, row["name"])
			if row["default"] == "true" {
				foundDefaultDomain = true
			}
		}

		if !foundDefaultDomain {
			t.Fatalf("Expected to find existing default domain, but did not: %#v", out)
		}
	})

	var newDomain string

	for domain, _ := range newDomainsSelection {
		newDomain = domain
		break
	}

	if len(newDomain) == 0 {
		t.Fatalf("Expected to select new domain")
	}

	logger.Section("Change default domain", func() {
		knctl.RunWithOpts([]string{"domain", "create", "-d", newDomain, "--default"}, RunOpts{NoNamespace: true})

		out, _ := knctl.RunWithOpts([]string{"domain", "list", "--json"}, RunOpts{NoNamespace: true})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) == 0 {
			t.Fatalf("Expected to see at least one domain, but did not: '%s'", out)
		}

		var foundDomain bool

		for _, row := range resp.Tables[0].Rows {
			if row["name"] == newDomain {
				foundDomain = true
				if row["default"] != "true" {
					t.Fatalf("Expected to see domain set to default, but did not: '%s'", row)
				}
			}
		}

		if !foundDomain {
			t.Fatalf("Expected to find domain in the list, but did not: '%s'", out)
		}
	})
}
