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

func TestIngresses(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}

	logger.Section("Checking if at least one ingress available", func() {
		out, _ := knctl.RunWithOpts([]string{"list", "ingresses", "--json"}, RunOpts{NoNamespace: true})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) == 0 {
			t.Fatalf("Expected to see at least one ingress, but did not: '%s'", out)
		}

		row := resp.Tables[0].Rows[0]

		if len(row["addresses"]) == 0 {
			t.Fatalf("Expected ingress to have at least one address, but was: %#v", row)
		}

		if len(row["ports"]) == 0 {
			t.Fatalf("Expected ingress to have one at least one port, but was: %#v", row)
		}
	})
}
