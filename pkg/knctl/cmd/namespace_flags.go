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

package cmd

import (
	"fmt"
	"os"

	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type NamespaceFlags struct {
	Name string
}

func (s *NamespaceFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	name := flagsFactory.NewNamespaceNameFlag(&s.Name)
	cmd.Flags().VarP(name, "namespace", "n", "Specified namespace (or default from kubeconfig)")
}

type NamespaceNameFlag struct {
	value       *string
	depsFactory DepsFactory
}

var _ pflag.Value = &NamespaceNameFlag{}
var _ cobrautil.ResolvableFlag = &NamespaceNameFlag{}

func NewNamespaceNameFlag(value *string, depsFactory DepsFactory) *NamespaceNameFlag {
	return &NamespaceNameFlag{value, depsFactory}
}

func (s *NamespaceNameFlag) Set(val string) error {
	*s.value = val
	return nil
}

func (s *NamespaceNameFlag) Type() string   { return "string" }
func (s *NamespaceNameFlag) String() string { return "" } // default for usage

func (s *NamespaceNameFlag) Resolve() error {
	value, err := s.resolveValue()
	if err != nil {
		return err
	}

	*s.value = value

	return nil
}

func (s *NamespaceNameFlag) resolveValue() (string, error) {
	if s.value != nil && len(*s.value) > 0 {
		return *s.value, nil
	}

	envVal := os.Getenv("KNCTL_NAMESPACE")
	if len(envVal) > 0 {
		return envVal, nil
	}

	configVal, err := s.depsFactory.DefaultNamespace()
	if err != nil {
		return configVal, nil
	}

	if len(configVal) > 0 {
		return configVal, nil
	}

	return "", fmt.Errorf("Expected to non-empty namespace name")
}
