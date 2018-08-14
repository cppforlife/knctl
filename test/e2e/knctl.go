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
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type Knctl struct {
	t         *testing.T
	namespace string
	l         Logger
}

type RunOpts struct {
	NoNamespace  bool
	AllowError   bool
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
	CancelCh     chan struct{}
	Redact       bool
}

func (k Knctl) Run(args []string) string {
	out, _ := k.RunWithOpts(args, RunOpts{})
	return out
}

func (k Knctl) RunWithOpts(args []string, opts RunOpts) (string, error) {
	k.l.Debugf("Running '%s'...\n", k.cmdDesc(args, opts))

	if !opts.NoNamespace {
		args = append(args, []string{"-n", k.namespace}...)
	}

	var stderr, stdout bytes.Buffer

	cmd := exec.Command("knctl", args...)

	if opts.StderrWriter != nil {
		cmd.Stderr = opts.StderrWriter
	} else {
		cmd.Stderr = &stderr
	}

	if opts.StdoutWriter != nil {
		cmd.Stdout = opts.StdoutWriter
	} else {
		cmd.Stdout = &stdout
	}

	if opts.CancelCh != nil {
		go func() {
			select {
			case <-opts.CancelCh:
				cmd.Process.Signal(os.Interrupt)
			}
		}()
	}

	err := cmd.Run()
	stdoutStr := stdout.String()

	if err != nil {
		err = fmt.Errorf("Execution error: stdout: '%s' stderr: '%s' error: '%s'", stdoutStr, stderr.String(), err)

		if !opts.AllowError {
			k.t.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args, opts), err)
		}
	}

	return stdoutStr, err
}

func (k Knctl) cmdDesc(args []string, opts RunOpts) string {
	if opts.Redact {
		return "knctl -redacted-"
	}
	return fmt.Sprintf("knctl %s", strings.Join(args, " "))
}
