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

package cmd

import (
	"os"
	"path/filepath"

	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type KubeconfigFlags struct {
	Path *KubeconfigPathFlag
}

func (f *KubeconfigFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	f.Path = NewKubeconfigPathFlag()
	cmd.PersistentFlags().Var(f.Path, "kubeconfig", "Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)")
}

type KubeconfigPathFlag struct {
	value string
}

var _ pflag.Value = &KubeconfigPathFlag{}
var _ cobrautil.ResolvableFlag = &KubeconfigPathFlag{}

func NewKubeconfigPathFlag() *KubeconfigPathFlag {
	return &KubeconfigPathFlag{}
}

func (s *KubeconfigPathFlag) Set(val string) error {
	s.value = val
	return nil
}

func (s *KubeconfigPathFlag) Type() string   { return "string" }
func (s *KubeconfigPathFlag) String() string { return "" } // default for usage

func (s *KubeconfigPathFlag) Value() (string, error) {
	err := s.Resolve()
	if err != nil {
		return "", err
	}

	return s.value, nil
}

func (s *KubeconfigPathFlag) Resolve() error {
	if len(s.value) > 0 {
		return nil
	}

	s.value = s.resolveValue()

	return nil
}

func (s *KubeconfigPathFlag) resolveValue() string {
	path := os.Getenv("KNCTL_KUBECONFIG")
	if len(path) > 0 {
		return path
	}

	path = os.Getenv("KUBECONFIG")
	if len(path) > 0 {
		return path
	}

	return filepath.Join(s.homeDir(), ".kube", "config")
}

func (*KubeconfigPathFlag) homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
