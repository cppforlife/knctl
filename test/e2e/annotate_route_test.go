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
	"fmt"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"gopkg.in/yaml.v2"
)

func TestAnnotateRoute(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}

	const (
		routeName1  = "test-annotate-route-name1"
		routeName2  = "test-annotate-route-name2"
		routeName3  = "test-annotate-route-name3"
		serviceName = "test-annotate-route-service-name"
	)

	routeNames := map[string]bool{routeName1: true, routeName2: true, routeName3: true}

	cleanUp := func() {
		for name, _ := range routeNames {
			knctl.RunWithOpts([]string{"route", "delete", "--route", name}, RunOpts{AllowError: true})
		}
	}

	logger.Section("Delete previous routes if exists", cleanUp)
	defer cleanUp()

	const (
		annotationKey           = "custom-key"
		annotationValue         = "custom-val"
		annotationCustomNameKey = "custom-name"
	)

	logger.Section("Create 3 routes", func() {
		for name, _ := range routeNames {
			knctl.Run([]string{"rollout", "--route", name, "--service-percentage", serviceName + "=100%"})
		}
	})

	logger.Section("Checking that there are no annotations", func() {
		out := knctl.Run([]string{"route", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) < 3 {
			t.Fatalf("Expected to see 3 routes in the list of routes, but did not: '%s'", out)
		}

		for name, _ := range routeNames {
			var found bool

			for _, row := range resp.Tables[0].Rows {
				if name != row["name"] {
					continue
				}

				found = true
				var anns map[string]interface{}

				err := yaml.Unmarshal([]byte(row["annotations"]), &anns)
				if err != nil {
					t.Fatalf("Expected YAML unmarshaling to succeed: '%s'", err)
				}

				if _, found := anns[annotationKey]; found {
					t.Fatalf("Did not expect to find annotation in '%#v'", anns)
				}
			}

			if !found {
				t.Fatalf("Expected to find route '%s' in out: %s", name, out)
			}
		}
	})

	logger.Section("Annotating routes", func() {
		for name, _ := range routeNames {
			ann1 := fmt.Sprintf("%s=%s", annotationKey, annotationValue)
			ann2 := fmt.Sprintf("%s=%s", annotationCustomNameKey, name)
			knctl.Run([]string{"route", "annotate", "--route", name, "-a", ann1, "-a", ann2})
		}
	})

	logger.Section("Checking that there are annotations", func() {
		out := knctl.Run([]string{"route", "list", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		for name, _ := range routeNames {
			var found bool

			for _, row := range resp.Tables[0].Rows {
				if name != row["name"] {
					continue
				}

				found = true
				var anns map[string]interface{}

				err := yaml.Unmarshal([]byte(row["annotations"]), &anns)
				if err != nil {
					t.Fatalf("Expected YAML unmarshaling to succeed: '%s'", err)
				}

				if anns[annotationKey] != annotationValue {
					t.Fatalf("Expected route to be annotated, but was not '%#v'", anns)
				}
				if anns[annotationCustomNameKey] != row["name"] {
					t.Fatalf("Expected route to be annotated with a second annotation, but was not '%#v'", anns)
				}
			}

			if !found {
				t.Fatalf("Expected to find route '%s' in out: %s", name, out)
			}
		}
	})
}
