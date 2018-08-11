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
	"bytes"
	"strings"
	"testing"
)

func TestCLIErrorsForFlagsBeforeExtraArgs(t *testing.T) {
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, Logger{}}

	var stderr bytes.Buffer

	_, err := knctl.RunWithOpts(
		[]string{"create", "namespace", "test-ns"},
		RunOpts{StderrWriter: &stderr, NoNamespace: true, AllowError: true},
	)

	if err == nil {
		t.Fatalf("Expected to receive error")
	}

	stderrStr := stderr.String()

	// Required flag error is more friendlier than command does not accept extra arg
	if !strings.Contains(stderrStr, `Error: required flag(s) "namespace" not set`) {
		t.Fatalf("Expected to find required flag error, but was '%s'", stderrStr)
	}
}

func TestCLIErrorsForCommandGroups(t *testing.T) {
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, Logger{}}

	// For commands with children commands it's friendlier ux
	// to ignore extra arguments and show available subcommands
	cmdsWithSubcmds := []string{"list", "create", "delete", "tag", "untag", "annotate"}

	for _, cmd := range cmdsWithSubcmds {
		var stderr bytes.Buffer

		_, err := knctl.RunWithOpts(
			[]string{cmd, "test-subcmd"},
			RunOpts{StderrWriter: &stderr, NoNamespace: true, AllowError: true},
		)

		if err == nil {
			t.Fatalf("[cmd %s] Expected to receive error", cmd)
		}

		stderrStr := stderr.String()

		if !strings.Contains(stderrStr, "Error: Use one of available subcommands") {
			t.Fatalf("[cmd %s] Expected to find invalid command error in '%s'", cmd, stderrStr)
		}
	}
}
