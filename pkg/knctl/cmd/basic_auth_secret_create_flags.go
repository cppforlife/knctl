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

	"github.com/spf13/cobra"
)

type BasicAuthSecretCreateFlags struct {
	GenerateNameFlags GenerateNameFlags

	Type     string
	URL      string
	Username string
	Password string

	DockerHub bool
	GCR       bool

	ForPulling bool
}

func (s *BasicAuthSecretCreateFlags) Set(cmd *cobra.Command) {
	s.GenerateNameFlags.Set(cmd)

	cmd.Flags().StringVar(&s.Type, "type", "", "Set type (example: docker, ssh)")
	cmd.Flags().StringVar(&s.URL, "url", "", "Set url (example: https://index.docker.io/v1/, https://github.com)")

	cmd.Flags().StringVarP(&s.Username, "username", "u", "", "Set username")
	cmd.MarkFlagRequired("username")

	defaultPassword := os.Getenv("KNCTL_BASIC_AUTH_SECRET_PASSWORD")
	cmd.Flags().StringVarP(&s.Password, "password", "p", defaultPassword, "Set password ($KNCTL_BASIC_AUTH_SECRET_PASSWORD)")
	if len(defaultPassword) == 0 {
		cmd.MarkFlagRequired("password")
	}

	cmd.Flags().BoolVar(&s.DockerHub, "docker-hub", false, "Preconfigure type and url for Docker Hub registry")
	cmd.Flags().BoolVar(&s.GCR, "gcr", false, "Preconfigure type and url for gcr.io registry")

	cmd.Flags().BoolVar(&s.ForPulling, "for-pulling", false, "Convert to pull secret ('kubernetes.io/dockerconfigjson' type)")
}

func (s *BasicAuthSecretCreateFlags) BackfillTypeAndURL() error {
	if s.GCR || s.DockerHub {
		if len(s.Type) != 0 || len(s.URL) != 0 {
			return fmt.Errorf("Expected to not specify --type or --url when preconfigured registry flags are used")
		}
	}

	switch {
	case s.DockerHub:
		s.Type = "docker"
		s.URL = "https://index.docker.io/v1/"

	case s.GCR:
		s.Type = "docker"
		s.URL = "https://gcr.io"

	default:
		if len(s.Type) == 0 || len(s.URL) == 0 {
			return fmt.Errorf("Expected --type and --url to be non-empty when preconfigured registry flags are not used")
		}
	}

	return nil
}
