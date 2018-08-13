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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListServicesOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	NamespaceFlags NamespaceFlags
}

func NewListServicesOptions(ui ui.UI, depsFactory DepsFactory) *ListServicesOptions {
	return &ListServicesOptions{ui: ui, depsFactory: depsFactory}
}

func NewListServicesCmd(o *ListServicesOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services",
		Aliases: servicesAliases,
		Short:   "List services",
		Long:    "List all services in a namespace",
		Example: `
  # List all services in namespace 'ns1'
  knctl list services -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd)
	return cmd
}

func (o *ListServicesOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	services, err := servingClient.ServingV1alpha1().Services(o.NamespaceFlags.Name).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	internalDomainHeader := uitable.NewHeader("Internal Domain")
	internalDomainHeader.Hidden = true

	table := uitable.Table{
		Title: fmt.Sprintf("Services in namespace '%s'", o.NamespaceFlags.Name),

		Content: "services",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Domain"),
			internalDomainHeader,
			uitable.NewHeader("Annotations"),
			uitable.NewHeader("Age"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, svc := range services.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(svc.Name),
			uitable.NewValueString(svc.Status.Domain),
			uitable.NewValueString(svc.Status.DomainInternal),
			NewAnnotationsValue(svc.Annotations),
			NewValueAge(svc.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}
