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

package service

import (
	"time"

	cmdbld "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/spf13/cobra"
)

type DeployFlags struct {
	GenerateNameFlags    cmdcore.GenerateNameFlags
	BuildCreateArgsFlags cmdbld.CreateArgsFlags

	Image string
	Env   []string

	WatchRevisionReady            bool
	WatchRevisionReadyMaxDuration time.Duration
	WatchPodLogs                  bool

	RemoveKnctlDeployEnvVar bool
}

func (s *DeployFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)
	s.BuildCreateArgsFlags.SetWithBuildPrefix(cmd, flagsFactory)

	// TODO separate service account for pulling?

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")

	cmd.Flags().BoolVar(&s.WatchRevisionReady, "watch-revision-ready", true, "Wait for new revision to become ready")
	cmd.Flags().DurationVar(&s.WatchRevisionReadyMaxDuration, "watch-revision-ready-max-duration",
		5*time.Minute, "Maximum duration of time to wait for new revision to become ready")
	cmd.Flags().BoolVar(&s.WatchPodLogs, "watch-pod-logs", true, "Watch pod logs for new revision")

	cmd.Flags().StringSliceVarP(&s.Env, "env", "e", nil, "Set environment variable (format: key=value) (can be specified multiple times)")
}
