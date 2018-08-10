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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type Knctl struct {
	t         *testing.T
	namespace string
	l         Logger
}

type RunOpts struct {
	NoNamespace bool
	AllowError  bool
	CancelCh    chan struct{}
}

func (k Knctl) Run(args []string) string {
	out, _ := k.RunWithOpts(args, RunOpts{})
	return out
}

func (k Knctl) RunWithOpts(args []string, opts RunOpts) (string, error) {
	k.l.Debugf("Running '%s'...\n", k.cmdDesc(args))

	if !opts.NoNamespace {
		args = append(args, []string{"-n", k.namespace}...)
	}

	var stderr bytes.Buffer

	cmd := exec.Command("knctl", args...)
	cmd.Stderr = &stderr

	if opts.CancelCh != nil {
		go func() {
			select {
			case <-opts.CancelCh:
				cmd.Process.Signal(os.Interrupt)
			}
		}()
	}

	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("Execution error: stderr: '%s' error: '%s'", stderr.String(), err)

		if !opts.AllowError {
			k.t.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args), err)
		}
	}

	return string(out), err
}

func (k Knctl) cmdDesc(args []string) string {
	return fmt.Sprintf("knctl %s", strings.Join(args, " "))
}

type Logger struct{}

func (l Logger) Section(msg string, f func()) {
	fmt.Printf("==> %s\n", msg)
	f()
}

func (l Logger) Debugf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

type Curl struct {
	t     *testing.T
	knctl Knctl
}

func (c Curl) WaitForContent(serviceName, expectedContent string) {
	var curledSuccessfully bool
	var out string

	for i := 0; i < 100; i++ {
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
