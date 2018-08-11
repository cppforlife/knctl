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

type CreateServiceAccountOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceAccountFlags       ServiceAccountFlags
	ServiceAccountCreateFlags ServiceAccountCreateFlags
}

func NewCreateServiceAccountOptions(ui ui.UI, depsFactory DepsFactory) *CreateServiceAccountOptions {
	return &CreateServiceAccountOptions{ui: ui, depsFactory: depsFactory}
}

func NewCreateServiceAccountCmd(o *CreateServiceAccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-account",
		Aliases: []string{"sa"},
		Short:   "Create service account",
		Long: `Create service account.

Use 'kubectl delete serviceaccount <name> -n <namespace>' to delete service account.`,
		Example: `
  # Create service account 'sa1' with two secrets in namespace 'ns1'
  knctl create service-account -a sa1 --secret secret1 --secret secret2 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceAccountFlags.Set(cmd)
	o.ServiceAccountCreateFlags.Set(cmd)
	return cmd
}

func (o *CreateServiceAccountOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.ServiceAccountFlags.Name, // TODO generate name
			Namespace: o.ServiceAccountFlags.NamespaceFlags.Name,
		},
	}

	for _, secretName := range o.ServiceAccountCreateFlags.Secrets {
		serviceAccount.Secrets = append(serviceAccount.Secrets, corev1.ObjectReference{Name: secretName})
	}

	_, err = coreClient.CoreV1().ServiceAccounts(o.ServiceAccountFlags.NamespaceFlags.Name).Create(serviceAccount)
	if err != nil {
		return fmt.Errorf("Creating service account: %s", err)
	}

	return nil
}
