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

package revision

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ServiceFlags cmdflags.ServiceFlags
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: cmdcore.ListAliases,
		Short:   "List revisions",
		Long:    "List all revisions for a service",
		Example: `
  # List all revisions for service 'svc1' in namespace 'ns1' 
  knctl revision list -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.SetOptional(cmd, flagsFactory)
	return cmd
}

func (o *ListOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tableTitle := "Revisions"
	serviceHeader := uitable.NewHeader("Service")
	listOpts := metav1.ListOptions{}

	if len(o.ServiceFlags.Name) > 0 {
		tableTitle += fmt.Sprintf(" for service '%s'", o.ServiceFlags.Name)
		serviceHeader.Hidden = true

		listOpts.LabelSelector = labels.Set(map[string]string{
			serving.ConfigurationLabelKey: o.ServiceFlags.Name,
		}).String()

		// Verify that service exists by this name
		_, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
	}

	revisions, err := servingClient.ServingV1alpha1().Revisions(o.ServiceFlags.NamespaceFlags.Name).List(listOpts)
	if err != nil {
		return err
	}

	routes, err := servingClient.ServingV1alpha1().Routes(o.ServiceFlags.NamespaceFlags.Name).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title:   tableTitle,
		Content: "revisions",

		Header: []uitable.Header{
			serviceHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Tags"),
			uitable.NewHeader("Annotations"),
			uitable.NewHeader("Conditions"),
			uitable.NewHeader("Age"),
			uitable.NewHeader("Traffic"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 5, Asc: false}, // Show latest first
		},
	}

	for _, rev := range revisions.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(rev.Labels[serving.ConfigurationLabelKey]),
			uitable.NewValueString(rev.Name),
			uitable.NewValueStrings(ctlservice.NewTags(servingClient).List(rev)),
			cmdcore.NewAnnotationsValue(rev.Annotations),
			cmdcore.NewConditionsValue(rev.Status.Conditions),
			cmdcore.NewValueAge(rev.CreationTimestamp.Time),
			NewTrafficValue(rev, routes.Items),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func NewTrafficValue(revision v1alpha1.Revision, routes []v1alpha1.Route) uitable.Value {
	var result []string
	for _, route := range routes {
		// Show based on actual configuration of the targets,
		// not based on desired configuration
		for _, target := range route.Status.Traffic {
			if target.RevisionName == revision.Name {
				domain := route.Status.Domain
				if len(target.Name) > 0 {
					domain += target.Name + "." + domain
				}
				result = append(result, fmt.Sprintf("%d%% -> %s", target.Percent, domain))
			}
		}
	}
	return uitable.NewValueStrings(result)
}
