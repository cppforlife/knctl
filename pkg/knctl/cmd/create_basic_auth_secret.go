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

type CreateBasicAuthSecretOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	SecretFlags                SecretFlags
	BasicAuthSecretCreateFlags BasicAuthSecretCreateFlags
}

func NewCreateBasicAuthSecretOptions(ui ui.UI, depsFactory DepsFactory) *CreateBasicAuthSecretOptions {
	return &CreateBasicAuthSecretOptions{ui: ui, depsFactory: depsFactory}
}

func NewCreateBasicAuthSecretCmd(o *CreateBasicAuthSecretOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basic-auth-secret",
		Short: "Create basic auth secret",
		Long: `Create basic auth secret.

Use 'kubectl delete secret <name> -n <namespace>' to delete secret.`,
		Example: `
  # Create SSH basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --type ssh --url github.com --username username --password password -n ns1

  # Create Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --docker-hub --username username --password password -n ns1

  # Create GCR.io registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --gcr --username username --password password -n ns1

  # Create generic Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --type docker --url https://registry.domain.com/ --username username --password password -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.SecretFlags.Set(cmd)
	o.BasicAuthSecretCreateFlags.Set(cmd)
	return cmd
}

func (o *CreateBasicAuthSecretOptions) Run() error {
	err := o.BasicAuthSecretCreateFlags.BackfillTypeAndURL()
	if err != nil {
		return err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.SecretFlags.Name, // TODO generate name
			Namespace: o.SecretFlags.NamespaceFlags.Name,
			Annotations: map[string]string{
				fmt.Sprintf("build.knative.dev/%s-0", o.BasicAuthSecretCreateFlags.Type): o.BasicAuthSecretCreateFlags.URL,
			},
		},
		Type: corev1.SecretTypeBasicAuth,
		StringData: map[string]string{
			corev1.BasicAuthUsernameKey: o.BasicAuthSecretCreateFlags.Username,
			corev1.BasicAuthPasswordKey: o.BasicAuthSecretCreateFlags.Password,
		},
	}

	_, err = coreClient.CoreV1().Secrets(o.SecretFlags.NamespaceFlags.Name).Create(secret)
	if err != nil {
		return fmt.Errorf("Creating basic auth secret: %s", err)
	}

	return nil
}
