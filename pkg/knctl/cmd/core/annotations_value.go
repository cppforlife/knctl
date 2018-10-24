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

package core

import (
	"strings"

	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/knative/serving/pkg/apis/serving"
)

const (
	kubectlLastAppliedAnnotation = "kubectl.kubernetes.io/last-applied-configuration"
)

func NewAnnotationsValue(anns map[string]string) uitable.Value {
	result := map[string]string{}
	for k, v := range anns {
		if !strings.HasPrefix(k, serving.GroupName) && k != kubectlLastAppliedAnnotation {
			result[k] = v
		}
	}
	return uitable.NewValueInterface(result)
}
