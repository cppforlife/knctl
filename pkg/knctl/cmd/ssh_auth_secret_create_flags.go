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

type SSHAuthSecretCreateFlags struct {
	GenerateNameFlags GenerateNameFlags

	Type       string
	URL        string
	PrivateKey string
	KnownHosts string

	Github bool
}

func (s *SSHAuthSecretCreateFlags) Set(cmd *cobra.Command) {
	s.GenerateNameFlags.Set(cmd)

	cmd.Flags().StringVar(&s.Type, "type", "", "Set type (example: git)")
	cmd.Flags().StringVar(&s.URL, "url", "", "Set url (example: github.com)")

	defaultKey := os.Getenv("KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY")
	cmd.Flags().StringVar(&s.PrivateKey, "private-key", defaultKey, "Set private key in PEM format ($KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY)")
	if len(defaultKey) == 0 {
		cmd.MarkFlagRequired("private-key")
	}

	cmd.Flags().StringVar(&s.KnownHosts, "known-hosts", "", "Set known hosts")

	cmd.Flags().BoolVar(&s.Github, "github", false, "Preconfigure type and url for Github.com Git access")
}

func (s *SSHAuthSecretCreateFlags) BackfillTypeAndURL() error {
	if s.Github {
		if len(s.Type) != 0 || len(s.URL) != 0 {
			return fmt.Errorf("Expected to not specify --type or --url when preconfigured flags are used")
		}
	}

	switch {
	case s.Github:
		s.Type = "git"
		s.URL = "github.com"

	default:
		if len(s.Type) == 0 || len(s.URL) == 0 {
			return fmt.Errorf("Expected --type and --url to be non-empty when preconfigured flags are not used")
		}
	}

	return nil
}
