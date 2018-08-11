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
	"strconv"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
)

type ListIngressesOptions struct {
	ui          ui.UI
	depsFactory DepsFactory
}

func NewListIngressesOptions(ui ui.UI, depsFactory DepsFactory) *ListIngressesOptions {
	return &ListIngressesOptions{ui: ui, depsFactory: depsFactory}
}

func NewListIngressesCmd(o *ListIngressesOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ingresses",
		Aliases: []string{"i", "is", "ing", "ings", "ingress"},
		Short:   "List ingresses",
		Long:    "List all ingresses labeled as `knative: ingressgateway` in Istio's namespace",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	return cmd
}

func (o *ListIngressesOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	services, err := IngressServices{coreClient}.List()
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: "Ingresses",

		Content: "ingresses",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Addresses"),
			uitable.NewHeader("Ports"),
			uitable.NewHeader("Age"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, svc := range services.Items {
		addrs := []string{}
		ports := []string{} // TODO int32

		for _, ing := range svc.Status.LoadBalancer.Ingress {
			if len(ing.IP) > 0 {
				addrs = append(addrs, ing.IP)
			}
			if len(ing.Hostname) > 0 {
				addrs = append(addrs, ing.Hostname)
			}
		}

		for _, port := range svc.Spec.Ports {
			ports = append(ports, strconv.Itoa(int(port.Port)))
		}

		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(svc.Name),
			uitable.NewValueStrings(addrs),
			uitable.NewValueStrings(ports),
			NewValueAge(svc.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}
