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
	"strconv"
	"time"

	cmdbld "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	"github.com/spf13/cobra"
)

type DeployFlags struct {
	GenerateNameFlags    cmdcore.GenerateNameFlags
	BuildCreateArgsFlags cmdbld.CreateArgsFlags
	TagFlags             cmdflags.TagFlags
	AnnotateFlags        cmdflags.AnnotateFlags

	Image         string
	EnvVars       []string
	EnvSecrets    []string
	EnvConfigMaps []string

	ContainerConcurrency *int
	MinScale             *int
	MaxScale             *int

	WatchRevisionReady        bool
	WatchRevisionReadyTimeout time.Duration

	WatchPodLogs             bool
	WatchPodLogsIndefinitely bool

	ManagedRoute bool

	RemoveKnctlDeployEnvVar bool
}

func (s *DeployFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)
	s.BuildCreateArgsFlags.SetWithBuildPrefix(cmd, flagsFactory)
	s.TagFlags.Set(cmd, flagsFactory)
	s.AnnotateFlags.Set(cmd, flagsFactory)

	// TODO separate service account for pulling?

	cmd.Flags().StringVarP(&s.Image, "image", "i", "", "Set image URL")
	cmd.MarkFlagRequired("image")

	cmd.Flags().BoolVar(&s.WatchRevisionReady, "watch-revision-ready", true, "Wait for new revision to become ready")
	cmd.Flags().DurationVar(&s.WatchRevisionReadyTimeout, "watch-revision-ready-timeout",
		5*time.Minute, "Set timeout for waiting for new revision to become ready")

	cmd.Flags().BoolVar(&s.WatchPodLogs, "watch-pod-logs", true, "Watch pod logs for new revision")
	cmd.Flags().BoolVarP(&s.WatchPodLogsIndefinitely, "watch-pod-logs-indefinitely", "l",
		false, "Watch pod logs for new revision indefinitely")

	cmd.Flags().StringArrayVarP(&s.EnvVars, "env", "e", nil, "Set environment variable (format: ENV_KEY=value) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.EnvSecrets, "env-secret", nil, "Set environment variable from a secret (format: ENV_KEY=secret-name/key) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.EnvConfigMaps, "env-config-map", nil, "Set environment variable from a config map (format: ENV_KEY=config-map-name/key) (can be specified multiple times)")

	cmd.Flags().Var(newDefaultlessIntValue(&s.ContainerConcurrency), "container-concurrency", "Set container concurrency")
	cmd.Flags().Var(newDefaultlessIntValue(&s.MinScale), "min-scale", "Set autoscaling rule for minimum number of containers")
	cmd.Flags().Var(newDefaultlessIntValue(&s.MaxScale), "max-scale", "Set autoscaling rule for maximum number of containers")

	cmd.Flags().BoolVar(&s.ManagedRoute, "managed-route", true, "Custom route configuration")
}

type defaultlessIntValue struct {
	val **int
}

func newDefaultlessIntValue(p **int) *defaultlessIntValue {
	return &defaultlessIntValue{p}
}

func (i *defaultlessIntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	val := int(v)
	(*i.val) = &val
	return nil
}

func (i *defaultlessIntValue) Type() string {
	return "int"
}

func (i *defaultlessIntValue) String() string {
	if i.val == nil || *i.val == nil {
		return "unspecified"
	}
	return strconv.Itoa(int(**i.val))
}
