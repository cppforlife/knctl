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
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type SSHAuthSecretCreateFlags struct {
	GenerateNameFlags GenerateNameFlags

	Type string
	URL  string

	PrivateKey     string
	PrivateKeyPath string

	KnownHosts string

	Github bool
}

func (s *SSHAuthSecretCreateFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	s.GenerateNameFlags.Set(cmd, flagsFactory)

	cmd.Flags().StringVar(&s.Type, "type", "", "Set type (example: git)")
	cmd.Flags().StringVar(&s.URL, "url", "", "Set url (example: github.com)")

	defaultKey := os.Getenv("KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY")
	cmd.Flags().StringVar(&s.PrivateKey, "private-key", defaultKey, "Set private key in PEM format ($KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY)")
	cmd.Flags().StringVar(&s.PrivateKeyPath, "private-key-path", "", "Set private key in PEM format from file path")

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

func (s *SSHAuthSecretCreateFlags) BackfillPrivateKey() error {
	if len(s.PrivateKey) > 0 && len(s.PrivateKeyPath) > 0 {
		return fmt.Errorf("Expected to not find --private-key and --private-key-path specified together")
	}

	if len(s.PrivateKey) == 0 && len(s.PrivateKeyPath) == 0 {
		return fmt.Errorf("Expected to find --private-key or --private-key-path specified")
	}

	if len(s.PrivateKeyPath) > 0 {
		contents, err := ioutil.ReadFile(s.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("Reading private key file path: %s", err)
		}

		s.PrivateKey = string(contents)
	}

	return nil
}

func (s *SSHAuthSecretCreateFlags) Validate() error {
	pemBlock, _ := pem.Decode([]byte(s.PrivateKey))
	if pemBlock == nil {
		var hint string
		if strings.Contains(s.PrivateKey, "\\n") {
			hint = " (it appears that newline characters were escaped)"
		}
		return fmt.Errorf("Expected to find at least one PEM block in private key%s", hint)
	}

	return nil
}
