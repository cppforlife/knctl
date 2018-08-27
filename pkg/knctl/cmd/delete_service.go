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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteServiceOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags ServiceFlags
}

func NewDeleteServiceOptions(ui ui.UI, depsFactory DepsFactory) *DeleteServiceOptions {
	return &DeleteServiceOptions{ui: ui, depsFactory: depsFactory}
}

func NewDeleteServiceCmd(o *DeleteServiceOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: deleteAliases,
		Short:   "Delete service",
		Example: `
  # Delete service 'svc1' in namespace 'ns1'
  knctl service delete -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *DeleteServiceOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	err = servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Delete(o.ServiceFlags.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Deleting service: %s", err)
	}

	// TODO idempotent?

	return nil
}
