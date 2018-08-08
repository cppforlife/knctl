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
	"github.com/spf13/cobra"
)

type DeployFlags struct {
	BuildCreateFlags BuildCreateFlags

	Image string
	Env   []string
}

func (s *DeployFlags) Set(cmd *cobra.Command) {
	s.BuildCreateFlags.Set(cmd, BuildCreateFlagsOpts{Optional: true, NoImage: true})

	// TODO separate service account for pulling?

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")

	cmd.Flags().StringSliceVarP(&s.Env, "env", "e", nil, "Set environment variable (format: key=value) (can be specified multiple times)")
}
