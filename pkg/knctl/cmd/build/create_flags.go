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

package build

import (
	"time"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/spf13/cobra"
)

type CreateFlags struct {
	GenerateNameFlags cmdcore.GenerateNameFlags
	CreateArgsFlags
}

type CreateArgsFlags struct {
	ctlbuild.BuildSpecOpts
}

func (s *CreateFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)
	s.CreateArgsFlags.Set(cmd, flagsFactory)

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")
}

func (s *CreateFlags) Validate() error {
	return s.CreateArgsFlags.Validate()
}

func (s *CreateArgsFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.setWithPrefix("", cmd, flagsFactory)
}

func (s *CreateArgsFlags) SetWithBuildPrefix(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.setWithPrefix("build", cmd, flagsFactory)
}

func (s *CreateArgsFlags) setWithPrefix(prefix string, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	if len(prefix) > 0 {
		prefix += "-"
	}

	cmd.Flags().StringVarP(&s.SourceDirectory, "directory", "d", "", "Set source code directory")

	cmd.Flags().StringVar(&s.GitURL, "git-url", "", "Set Git URL")
	cmd.Flags().StringVar(&s.GitRevision, "git-revision", "", "Set Git revision (examples: https://git-scm.com/docs/gitrevisions#_specifying_revisions)")

	cmd.Flags().StringVar(&s.ServiceAccountName, "service-account", "", "Set service account name for building") // TODO separate

	cmd.Flags().StringVar(&s.Template, "template", "", "Set template name")
	cmd.Flags().StringSliceVar(&s.TemplateArgs, "template-arg", nil, "Set template argument (format: key=value) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.TemplateEnv, "template-env", nil, "Set template environment variable (format: key=value) (can be specified multiple times)")

	cmd.Flags().DurationVar(&s.Timeout, prefix+"timeout", time.Duration(0), "Set timeout for building stage (Knative Build has a 10m default)")
}

func (s *CreateArgsFlags) IsProvided() bool {
	return len(s.SourceDirectory) > 0 || len(s.GitURL) > 0
}

func (s *CreateArgsFlags) Validate() error {
	return nil // TODO better error messages?
}
