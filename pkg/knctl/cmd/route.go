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
	"strconv"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RouteOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	RouteFlags   RouteFlags
	TrafficFlags TrafficFlags
}

func NewRouteOptions(ui ui.UI, depsFactory DepsFactory) *RouteOptions {
	return &RouteOptions{ui: ui, depsFactory: depsFactory}
}

func NewRouteCmd(o *RouteOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Configure route",
		Example: `
  # Set traffic percentages for service 'svc1' in namespace 'ns1'
  knctl route --route rt1 -p svc1:latest=20% -p svc1:previous=80% -n ns1

  # Roll back traffic for previous revision of service 'svc1' in namespace 'ns1'
  knctl route --route rt1 -p svc1:previous=100% -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RouteFlags.Set(cmd, flagsFactory)
	o.TrafficFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *RouteOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	route, err := servingClient.ServingV1alpha1().Routes(o.RouteFlags.NamespaceFlags.Name).Get(o.RouteFlags.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Getting route: %s", err)
		}

		route = &v1alpha1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      o.RouteFlags.Name, // TODO generate name
				Namespace: o.RouteFlags.NamespaceFlags.Name,
			},
		}
	}

	var targets []v1alpha1.TrafficTarget

	for _, traffic := range o.TrafficFlags.Percentages {
		pieces := strings.SplitN(traffic, "=", 2)
		if len(pieces) != 2 {
			return fmt.Errorf("Expected percentage to be in format 'revision=percentage'")
		}

		percent, err := strconv.Atoi(strings.TrimSuffix(pieces[1], "%"))
		if err != nil {
			return fmt.Errorf("Expected percentage value to be an integer")
		}

		if percent < 0 || percent > 100 {
			return fmt.Errorf("Expected percentage value to be between 0%% and 100%%")
		}

		revFlags := RevisionFlags{Name: pieces[0], NamespaceFlags: o.RouteFlags.NamespaceFlags}

		revision, err := NewRevisionReference(revFlags, tags, servingClient).Revision()
		if err != nil {
			return err
		}

		targets = append(targets, v1alpha1.TrafficTarget{
			RevisionName: revision.Name,
			Percent:      percent,
			// TODO ConfiguratioName 'service:'?
			// TODO Name
		})
	}

	route.Spec.Traffic = targets

	return o.createOrUpdate(servingClient, route)
}

func (o *RouteOptions) createOrUpdate(servingClient servingclientset.Interface, route *v1alpha1.Route) error {
	if len(route.ResourceVersion) > 0 {
		var updateErr error

		// TODO better retry functionality
		for i := 0; i < 5; i++ {
			_, updateErr = servingClient.ServingV1alpha1().Routes(o.RouteFlags.NamespaceFlags.Name).Update(route)
			if updateErr == nil {
				return nil
			}
		}

		return fmt.Errorf("Updating route: %s", updateErr)
	}

	_, createErr := servingClient.ServingV1alpha1().Routes(o.RouteFlags.NamespaceFlags.Name).Create(route)
	if createErr != nil {
		return fmt.Errorf("Creating route: %s", createErr)
	}

	return nil
}
