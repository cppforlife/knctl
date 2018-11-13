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

package route

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ShowOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	RouteFlags RouteFlags
}

func NewShowOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *ShowOptions {
	return &ShowOptions{ui: ui, depsFactory: depsFactory}
}

func NewShowCmd(o *ShowOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show route",
		Long:  "Show route details in a namespace",
		Example: `
  # Show details for route 'route1' in namespace 'ns1'
  knctl route show --route route1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RouteFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *ShowOptions) Run() error {
	routeClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	route, err := routeClient.ServingV1alpha1().Routes(o.RouteFlags.NamespaceFlags.Name).Get(o.RouteFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	o.printStatus(route)
	o.printTargets(route)

	cmdcore.NewConditionsTable(route.Status.Conditions).Print(o.ui)

	return nil
}

func (o *ShowOptions) printStatus(route *v1alpha1.Route) {
	table := uitable.Table{
		Title: fmt.Sprintf("Route '%s'", o.RouteFlags.Name),
		// TODO Content: "route",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Domain"),
			uitable.NewHeader("Internal Domain"),
			uitable.NewHeader("Age"),
		},

		Transpose: true,
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(route.Name),
		uitable.NewValueString(route.Status.Domain),
		uitable.NewValueString(route.Status.DomainInternal),
		cmdcore.NewValueAge(route.CreationTimestamp.Time),
	})

	o.ui.PrintTable(table)
}

func (o *ShowOptions) printTargets(route *v1alpha1.Route) {
	table := uitable.Table{
		Title: "Targets",
		// TODO Content: "targets",

		Header: []uitable.Header{
			uitable.NewHeader("Percent"),
			uitable.NewHeader("Revision"),
			uitable.NewHeader("Service"),
			uitable.NewHeader("Domain"),
		},
	}

	for _, tr := range route.Status.Traffic {
		domain := route.Status.Domain
		if len(tr.Name) > 0 {
			domain = tr.Name + "." + route.Status.Domain
		}

		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueSuffix(uitable.NewValueInt(tr.Percent), "%"),
			uitable.NewValueString(tr.RevisionName),
			uitable.NewValueString(tr.ConfigurationName),
			uitable.NewValueString(domain),
		})
	}

	o.ui.PrintTable(table)
}
