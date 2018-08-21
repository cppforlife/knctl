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

package util

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// Retry is different from wait.Poll because
// it does not stop retrying when error is encountered
func Retry(interval, timeout time.Duration, condFunc wait.ConditionFunc) error {
	var lastErr error
	var times int

	wait.Poll(interval, timeout, func() (bool, error) {
		done, err := condFunc()
		lastErr = err
		times++
		return done, nil
	})

	if lastErr != nil {
		return fmt.Errorf("Retried %d times: %s", times, lastErr)
	}

	return nil
}
