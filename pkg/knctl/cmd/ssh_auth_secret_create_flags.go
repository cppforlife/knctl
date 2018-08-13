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
	"os"

	"github.com/spf13/cobra"
)

type SSHAuthSecretCreateFlags struct {
	GenerateNameFlags GenerateNameFlags

	URL        string
	PrivateKey string
	KnownHosts string
}

func (s *SSHAuthSecretCreateFlags) Set(cmd *cobra.Command) {
	s.GenerateNameFlags.Set(cmd)

	cmd.Flags().StringVar(&s.URL, "url", "", "Set url (example: github.com)")
	cmd.MarkFlagRequired("url")

	defaultKey := os.Getenv("KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY")
	cmd.Flags().StringVar(&s.PrivateKey, "private-key", defaultKey, "Set private key in PEM format ($KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY)")
	if len(defaultKey) == 0 {
		cmd.MarkFlagRequired("private-key")
	}

	cmd.Flags().StringVar(&s.KnownHosts, "known-hosts", "", "Set known hosts")
}
