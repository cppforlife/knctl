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
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateServiceAccountOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceAccountFlags       ServiceAccountFlags
	ServiceAccountCreateFlags ServiceAccountCreateFlags
}

func NewCreateServiceAccountOptions(ui ui.UI, depsFactory DepsFactory) *CreateServiceAccountOptions {
	return &CreateServiceAccountOptions{ui: ui, depsFactory: depsFactory}
}

func NewCreateServiceAccountCmd(o *CreateServiceAccountOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create service account",
		Long: `Create service account.

Use 'kubectl get serviceaccount -n <namespace>' to list service accounts.
Use 'kubectl delete serviceaccount <name> -n <namespace>' to delete service account.`,
		Example: `
  # Create service account 'sa1' with two secrets in namespace 'ns1'
  knctl service-account create -a sa1 --secret secret1 --secret secret2 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceAccountFlags.Set(cmd, flagsFactory)
	o.ServiceAccountCreateFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *CreateServiceAccountOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: o.ServiceAccountCreateFlags.GenerateNameFlags.Apply(metav1.ObjectMeta{
			Name:      o.ServiceAccountFlags.Name,
			Namespace: o.ServiceAccountFlags.NamespaceFlags.Name,
		}),
	}

	o.populateSecrets(serviceAccount, coreClient)

	createdServiceAccount, err := coreClient.CoreV1().ServiceAccounts(o.ServiceAccountFlags.NamespaceFlags.Name).Create(serviceAccount)
	if err != nil {
		return fmt.Errorf("Creating service account: %s", err)
	}

	o.printTable(createdServiceAccount)

	return nil
}

func (o *CreateServiceAccountOptions) populateSecrets(sa *corev1.ServiceAccount, coreClient kubernetes.Interface) error {
	for _, secretName := range o.ServiceAccountCreateFlags.Secrets {
		secret, err := coreClient.CoreV1().Secrets(o.ServiceAccountFlags.NamespaceFlags.Name).Get(secretName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Getting secret '%s': %s", secretName, err)
		}

		if secret.Type == corev1.SecretTypeDockerConfigJson {
			sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{Name: secretName})
		} else {
			sa.Secrets = append(sa.Secrets, corev1.ObjectReference{Name: secretName})
		}
	}

	// Explicit image pull secrets
	for _, secretName := range o.ServiceAccountCreateFlags.ImagePullSecrets {
		sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{Name: secretName})
	}

	return nil
}

func (o *CreateServiceAccountOptions) printTable(sa *corev1.ServiceAccount) {
	table := uitable.Table{
		Header: []uitable.Header{
			uitable.NewHeader("Name"),
		},

		Transpose: true,

		Rows: [][]uitable.Value{
			{uitable.NewValueString(sa.Name)},
		},
	}

	o.ui.PrintTable(table)
}
