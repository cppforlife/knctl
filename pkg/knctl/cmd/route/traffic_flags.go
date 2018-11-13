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

package route

import (
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/spf13/cobra"
)

type TrafficFlags struct {
	GenerateNameFlags cmdcore.GenerateNameFlags

	RevisionPercentages []string
	ServicePercentages  []string
}

func (s *TrafficFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)

	cmd.Flags().StringSliceVarP(&s.RevisionPercentages, "percentage", "p", nil, "Set revision percentage (format: revision=percentage, example: app-00001=100%, app:latest=100%) (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&s.ServicePercentages, "service-percentage", nil, "Set service percentage (format: service=percentage, example: app=100%) (can be specified multiple times)")
}
