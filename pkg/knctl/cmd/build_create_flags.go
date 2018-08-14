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
	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	"github.com/spf13/cobra"
)

type BuildCreateFlags struct {
	GenerateNameFlags GenerateNameFlags
	BuildCreateArgsFlags
}

type BuildCreateArgsFlags struct {
	ctlbuild.BuildSpecOpts
}

func (s *BuildCreateFlags) Set(cmd *cobra.Command) {
	s.GenerateNameFlags.Set(cmd)
	s.BuildCreateArgsFlags.Set(cmd)

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")
}

func (s *BuildCreateFlags) Validate() error {
	return s.BuildCreateArgsFlags.Validate()
}

func (s *BuildCreateArgsFlags) Set(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&s.SourceDirectory, "directory", "d", "", "Set source code directory")

	cmd.Flags().StringVar(&s.GitURL, "git-url", "", "Set Git URL")
	cmd.Flags().StringVar(&s.GitRevision, "git-revision", "", "Set Git revision (examples: https://git-scm.com/docs/gitrevisions#_specifying_revisions)")

	cmd.Flags().StringVar(&s.ServiceAccountName, "service-account", "", "Set service account name for building") // TODO separate

	cmd.Flags().StringVar(&s.Template, "template", "", "Set template name")
	cmd.Flags().StringSliceVar(&s.TemplateArgs, "template-arg", nil, "Set template argument (format: key=value) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.TemplateEnv, "template-env", nil, "Set template environment variable (format: key=value) (can be specified multiple times)")
}

func (s *BuildCreateArgsFlags) IsProvided() bool {
	return len(s.SourceDirectory) > 0 || len(s.GitURL) > 0
}

func (s *BuildCreateArgsFlags) Validate() error {
	return nil // TODO better error messages?
}
