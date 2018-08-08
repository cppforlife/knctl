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
	t *testing.T
	l Logger
}

func (k Knctl) Run(args []string) string {
	out, err := k.RunWithErr(args)
	if err != nil {
		k.t.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args), err)
	}

	return out
}

func (k Knctl) RunWithErr(args []string) (string, error) {
	k.l.Debugf("Running '%s'...\n", k.cmdDesc(args))

	var stderr bytes.Buffer

	cmd := exec.Command("knctl", args...)
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("Execution error: stderr: '%s' error: '%s'", stderr.String(), err)
	}

	return string(out), err
}

func (k Knctl) RunWithCancel(args []string, cancelCh chan struct{}) (string, error) {
	k.l.Debugf("Running '%s'...\n", k.cmdDesc(args))

	var stderr bytes.Buffer

	cmd := exec.Command("knctl", args...)
	cmd.Stderr = &stderr

	go func() {
		select {
		case <-cancelCh:
			cmd.Process.Signal(os.Interrupt)
		}
	}()

	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("Execution error: stderr: '%s' error: '%s'", stderr.String(), err)
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
		out, _ = c.knctl.RunWithErr([]string{"curl", "-n", "default", "-s", serviceName})
		if strings.Contains(out, expectedContent) {
			curledSuccessfully = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !curledSuccessfully {
		c.t.Fatalf("Expected to see sample service output, but was: '%s'", out)
	}
}
