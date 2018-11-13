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
)

type Curl struct {
	t     *testing.T
	knctl Knctl
}

func (c Curl) WaitForContent(serviceName, expectedContent string) {
	var curledSuccessfully bool
	var out string

	for i := 0; i < 300; i++ {
		out, _ = c.knctl.RunWithOpts([]string{"curl", "-n", "default", "-s", serviceName}, RunOpts{AllowError: true})
		if strings.Contains(out, expectedContent) {
			curledSuccessfully = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !curledSuccessfully {
		c.t.Fatalf("Expected to find output '%s' in '%s' but did not", expectedContent, out)
	}
}

func (c Curl) WaitForRouteContent(routeName, expectedContent string) {
	var curledSuccessfully bool
	var out string

	for i := 0; i < 300; i++ {
		out, _ = c.knctl.RunWithOpts([]string{"route", "curl", "-n", "default", "--route", routeName}, RunOpts{AllowError: true})
		if strings.Contains(out, expectedContent) {
			curledSuccessfully = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !curledSuccessfully {
		c.t.Fatalf("Expected to find output '%s' in '%s' but did not", expectedContent, out)
	}
}
