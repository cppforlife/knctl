/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Open 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flags

import (
	"fmt"
	"strings"

	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/spf13/cobra"
)

type AnnotateFlags struct {
	Annotations []string
}

func (s *AnnotateFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	cmd.Flags().StringSliceVarP(&s.Annotations, "annotation", "a", nil, "Set annotation (format: key=value) (can be specified multiple times)")
}

func (s *AnnotateFlags) AsMap() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, kv := range s.Annotations {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected annotation to be in format 'KEY=VALUE'")
		}
		result[pieces[0]] = pieces[1]
	}

	return result, nil
}
