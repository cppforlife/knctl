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
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestLogsFollow(t *testing.T) {
	logger := Logger{}
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, logger}
	curl := Curl{t, knctl}

	const (
		serviceName         = "test-logs-follow-service-name"
		expectedContentRev1 = "TestRevisions_ContentRev1"
		expectedContentRev2 = "TestRevisions_ContentRev2"
		expectedContentRev3 = "TestRevisions_ContentRev3"
		expectedContentRev4 = "TestRevisions_ContentRev4"
	)

	logger.Section("Sleeping...", func() {
		// TODO otherwise 'no upstream healty' error happens
		// somehow caused by previous deploy in other tests
		time.Sleep(20 * time.Second)
	})

	logger.Section("Delete previous service with the same name if exists", func() {
		knctl.RunWithOpts([]string{"delete", "service", "-s", serviceName}, RunOpts{AllowError: true})
	})

	defer func() {
		knctl.RunWithOpts([]string{"delete", "service", "-s", serviceName}, RunOpts{AllowError: true})
	}()

	logger.Section("Deploy revision 1", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev1,
		})

		curl.WaitForContent(serviceName, expectedContentRev1)
	})

	cancelCh := make(chan struct{})
	doneCh := make(chan struct{})
	collectedLogs := &LogsWriter{}

	// Start tailing logs in the backgroud
	go func() {
		knctl.RunWithOpts([]string{"logs", "-s", serviceName, "-f"}, RunOpts{StdoutWriter: collectedLogs, CancelCh: cancelCh})
		doneCh <- struct{}{}
	}()

	logger.Section("Deploy revision 2 and check its logs", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev2,
		})

		curl.WaitForContent(serviceName, expectedContentRev2)
	})

	logger.Section("Deploy revision 3 and check its logs", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev3,
		})

		curl.WaitForContent(serviceName, expectedContentRev3)
	})

	logger.Section("Deploy revision 4 and check its logs", func() {
		knctl.Run([]string{
			"deploy",
			"-s", serviceName,
			"-i", "gcr.io/knative-samples/helloworld-go",
			"-e", "TARGET=" + expectedContentRev4,
		})

		curl.WaitForContent(serviceName, expectedContentRev4)
	})

	logger.Section("Check logs of service to make sure it includes logs from 4 revisions", func() {
		out := knctl.Run([]string{"list", "revisions", "-s", serviceName, "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))

		if len(resp.Tables[0].Rows) != 4 {
			t.Fatalf("Expected to see one revision in the list of revisions, but did not: '%s'", out)
		}

		checkLogs := func(currentLogs string, revisionRows []map[string]string) error {
			currentLogsLines := strings.Split(currentLogs, "\n")

			expectedLogLines := []string{
				"Hello world sample started.",
				"Hello world received a request.",
			}

			var matchedLines int

			for _, row := range revisionRows {
				expectedRevision := row["name"]

				for _, expectedLogLine := range expectedLogLines {
					var found bool
					for _, line := range currentLogsLines {
						if strings.HasPrefix(line, expectedRevision+" >") && strings.HasSuffix(line, expectedLogLine) {
							found = true
							matchedLines++
							break
						}
					}
					if !found {
						return fmt.Errorf("Expected to find log line '%s' for revision '%s' in service logs: '%s'", expectedLogLine, expectedRevision, currentLogs)
					}
				}
			}

			if matchedLines == 0 {
				return fmt.Errorf("Expected to have matched several lines")
			}

			return nil
		}

		var lastErr error

		// Try multiple times for all logs to make it
		for i := 0; i < 20; i++ {
			logger.Debugf("Trying to check logs try=%d\n", i)
			lastErr = checkLogs(collectedLogs.Current(), resp.Tables[0].Rows)
			if lastErr == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}

		if lastErr != nil {
			t.Fatalf(lastErr.Error())
		}
	})

	cancelCh <- struct{}{}
	<-doneCh

	logger.Section("Deleting service", func() {
		knctl.Run([]string{"delete", "service", "-s", serviceName})

		out := knctl.Run([]string{"list", "services", "--json"})
		if strings.Contains(out, serviceName) {
			t.Fatalf("Expected to not see sample service in the list of services, but was: %s", out)
		}
	})
}

type LogsWriter struct {
	lock   sync.RWMutex
	output []byte
}

var _ io.Writer = &LogsWriter{}

func (w *LogsWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.output = append(w.output, p...)
	return len(p), nil
}

func (w *LogsWriter) Current() string {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return string(w.output)
}
