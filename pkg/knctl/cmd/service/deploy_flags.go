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

	Image         string
	EnvVars       []string
	EnvSecrets    []string
	EnvConfigMaps []string

	WatchRevisionReady        bool
	WatchRevisionReadyTimeout time.Duration

	WatchPodLogs             bool
	WatchPodLogsIndefinitely bool

	CustomRoute bool

	RemoveKnctlDeployEnvVar bool
}

func (s *DeployFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)
	s.BuildCreateArgsFlags.SetWithBuildPrefix(cmd, flagsFactory)

	// TODO separate service account for pulling?

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")

	cmd.Flags().BoolVar(&s.WatchRevisionReady, "watch-revision-ready", true, "Wait for new revision to become ready")
	cmd.Flags().DurationVar(&s.WatchRevisionReadyTimeout, "watch-revision-ready-timeout",
		5*time.Minute, "Set timeout for waiting for new revision to become ready")

	cmd.Flags().BoolVar(&s.WatchPodLogs, "watch-pod-logs", true, "Watch pod logs for new revision")
	cmd.Flags().BoolVarP(&s.WatchPodLogsIndefinitely, "watch-pod-logs-indefinitely", "l",
		false, "Watch pod logs for new revision indefinitely")

	cmd.Flags().StringSliceVarP(&s.EnvVars, "env", "e", nil, "Set environment variable (format: ENV_KEY=value) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.EnvSecrets, "env-secret", nil, "Set environment variable from a secret (format: ENV_KEY=secret-name/key) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.EnvConfigMaps, "env-config-map", nil, "Set environment variable from a config map (format: ENV_KEY=config-map-name/key) (can be specified multiple times)")

	cmd.Flags().BoolVar(&s.CustomRoute, "custom-route", false, "Custom route configuration")
}
