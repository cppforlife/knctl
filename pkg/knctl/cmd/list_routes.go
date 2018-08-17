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
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListRoutesOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	NamespaceFlags NamespaceFlags
}

func NewListRoutesOptions(ui ui.UI, depsFactory DepsFactory) *ListRoutesOptions {
	return &ListRoutesOptions{ui: ui, depsFactory: depsFactory}
}

func NewListRoutesCmd(o *ListRoutesOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "routes",
		Aliases: []string{"route"},
		Short:   "List routes",
		Long:    "List all routes in a namespace",
		Example: `
  # List all routes in namespace 'ns1'
  knctl list routes -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *ListRoutesOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	routes, err := servingClient.ServingV1alpha1().Routes(o.NamespaceFlags.Name).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	internalDomainHeader := uitable.NewHeader("Internal Domain")
	internalDomainHeader.Hidden = true

	table := uitable.Table{
		Title: fmt.Sprintf("Routes in namespace '%s'", o.NamespaceFlags.Name),

		Content: "routes",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Traffic"),
			uitable.NewHeader("All Traffic Assigned"),
			uitable.NewHeader("Ready"),
			uitable.NewHeader("Domain"),
			internalDomainHeader,
			uitable.NewHeader("Age"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, route := range routes.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(route.Name),
			o.configurationValue(route),
			o.allTrafficAssignedValue(route),
			o.readyValue(route),
			uitable.NewValueString(route.Status.Domain),
			uitable.NewValueString(route.Status.DomainInternal),
			NewValueAge(route.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (*ListRoutesOptions) configurationValue(route v1alpha1.Route) uitable.ValueStrings {
	var dsts []string
	for _, target := range route.Spec.Traffic {
		dsts = append(dsts, fmt.Sprintf("%3d%% -> %s:%s", target.Percent, target.ConfigurationName, target.RevisionName))
	}
	return uitable.NewValueStrings(dsts)
}

func (*ListRoutesOptions) allTrafficAssignedValue(route v1alpha1.Route) ValueUnknownBool {
	cond := route.Status.GetCondition(v1alpha1.RouteConditionAllTrafficAssigned)
	if cond != nil {
		switch cond.Status {
		case corev1.ConditionTrue:
			result := true
			return NewValueUnknownBool(&result)
		case corev1.ConditionFalse:
			result := false
			return NewValueUnknownBool(&result)
		}
	}

	return NewValueUnknownBool(nil)
}

func (*ListRoutesOptions) readyValue(route v1alpha1.Route) ValueUnknownBool {
	cond := route.Status.GetCondition(v1alpha1.RouteConditionReady)
	if cond != nil {
		switch cond.Status {
		case corev1.ConditionTrue:
			result := true
			return NewValueUnknownBool(&result)
		case corev1.ConditionFalse:
			result := false
			return NewValueUnknownBool(&result)
		}
	}

	return NewValueUnknownBool(nil)
}
