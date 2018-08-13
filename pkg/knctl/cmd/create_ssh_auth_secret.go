/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
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

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateSSHAuthSecretOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	SecretFlags              SecretFlags
	SSHAuthSecretCreateFlags SSHAuthSecretCreateFlags
}

func NewCreateSSHAuthSecretOptions(ui ui.UI, depsFactory DepsFactory) *CreateSSHAuthSecretOptions {
	return &CreateSSHAuthSecretOptions{ui: ui, depsFactory: depsFactory}
}

func NewCreateSSHAuthSecretCmd(o *CreateSSHAuthSecretOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh-auth-secret",
		Aliases: []string{"sas"},
		Short:   "Create SSH auth secret",
		Long: `Create SSH auth secret.

Use 'kubectl delete secret <name> -n <namespace>' to delete secret.`,
		Example: `
  # Create SSH secret 'secret1' in namespace 'ns1'
  knctl create ssh-auth-secret -s secret1 --url github.com --private-key ... --known-hosts ... -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.SecretFlags.Set(cmd)
	o.SSHAuthSecretCreateFlags.Set(cmd)
	return cmd
}

func (o *CreateSSHAuthSecretOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.SecretFlags.Name, // TODO generate name
			Namespace: o.SecretFlags.NamespaceFlags.Name,
			Annotations: map[string]string{
				"build.knative.dev/git-0": o.SSHAuthSecretCreateFlags.URL,
			},
		},
		Type: corev1.SecretTypeSSHAuth,
		StringData: map[string]string{
			corev1.SSHAuthPrivateKey: o.SSHAuthSecretCreateFlags.PrivateKey,
		},
	}

	if len(o.SSHAuthSecretCreateFlags.KnownHosts) > 0 {
		secret.StringData["known_hosts"] = o.SSHAuthSecretCreateFlags.KnownHosts
	}

	_, err = coreClient.CoreV1().Secrets(o.SecretFlags.NamespaceFlags.Name).Create(secret)
	if err != nil {
		return fmt.Errorf("Creating ssh auth secret: %s", err)
	}

	return nil
}
