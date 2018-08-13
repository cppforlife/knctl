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
)

type CreateNamespaceOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	NamespaceFlags    NamespaceFlags
	GenerateNameFlags GenerateNameFlags
}

func NewCreateNamespaceOptions(ui ui.UI, depsFactory DepsFactory) *CreateNamespaceOptions {
	return &CreateNamespaceOptions{ui: ui, depsFactory: depsFactory}
}

func NewCreateNamespaceCmd(o *CreateNamespaceOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"ns"},
		Short:   "Create namespace",
		Long: `Create namespace.

Use 'kubectl delete ns <name>' to delete namespace.`,
		Example: `
  # Create namespace 'ns1'
  knctl create namespace -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd)
	o.GenerateNameFlags.Set(cmd)
	return cmd
}

func (o *CreateNamespaceOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	namespace := &corev1.Namespace{
		ObjectMeta: o.GenerateNameFlags.Apply(metav1.ObjectMeta{
			Name: o.NamespaceFlags.Name,
		}),
	}

	createdNamespace, err := coreClient.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return fmt.Errorf("Creating namespace: %s", err)
	}

	o.printTable(createdNamespace)

	// TODO idempotent?

	err = NewIstio(coreClient).EnableNamespace(o.NamespaceFlags.Name)
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateNamespaceOptions) printTable(ns *corev1.Namespace) {
	table := uitable.Table{
		Header: []uitable.Header{
			uitable.NewHeader("Name"),
		},

		Transpose: true,

		Rows: [][]uitable.Value{
			{uitable.NewValueString(ns.Name)},
		},
	}

	o.ui.PrintTable(table)
}
