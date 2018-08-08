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
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ListRevisionsOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags ServiceFlags
}

func NewListRevisionsOptions(ui ui.UI, depsFactory DepsFactory) *ListRevisionsOptions {
	return &ListRevisionsOptions{ui: ui, depsFactory: depsFactory}
}

func NewListRevisionsCmd(o *ListRevisionsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revisions",
		Aliases: []string{"r", "rs", "rev", "revs", "revision"},
		Short:   "List revisions",
		Long:    "List all revisions for a service",
		Example: `
  # List all revisions for service 'svc1' in namespace 'ns1' 
  knctl list revisions -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd)
	return cmd
}

func (o *ListRevisionsOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	listOpts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			serving.ConfigurationLabelKey: o.ServiceFlags.Name,
		}).String(),
	}

	revisions, err := servingClient.ServingV1alpha1().Revisions(o.ServiceFlags.NamespaceFlags.Name).List(listOpts)
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title:   fmt.Sprintf("Revisions for service '%s'", service.Name),
		Content: "revisions",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Tags"),
			uitable.NewHeader("Allocated Traffic %"),
			uitable.NewHeader("Serving State"),
			uitable.NewHeader("Annotations"),
			uitable.NewHeader("Created At"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 5, Asc: false}, // Show latest first
		},
	}

	for _, rev := range revisions.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(rev.Name),
			uitable.NewValueStrings(ctlservice.NewTags(servingClient).List(rev)),
			NewAllocatedTrafficPercentValue(service, rev),
			uitable.NewValueString(string(rev.Spec.ServingState)),
			NewAnnotationsValue(rev.Annotations),
			uitable.NewValueTime(rev.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func NewAllocatedTrafficPercentValue(svc *v1alpha1.Service, rev v1alpha1.Revision) uitable.Value {
	percent := 0
	for _, item := range svc.Status.Traffic {
		if item.RevisionName == rev.Name {
			percent = item.Percent
			break
		}
	}
	return uitable.NewValueSuffix(uitable.NewValueInt(percent), "%")
}

func NewAnnotationsValue(anns map[string]string) uitable.Value {
	result := map[string]string{}
	for k, v := range anns {
		if !strings.HasPrefix(k, serving.GroupName) {
			result[k] = v
		}
	}
	return uitable.NewValueInterface(result)
}
