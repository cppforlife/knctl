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

type DeleteBuildOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	BuildFlags BuildFlags
}

func NewDeleteBuildOptions(ui ui.UI, depsFactory DepsFactory) *DeleteBuildOptions {
	return &DeleteBuildOptions{ui: ui, depsFactory: depsFactory}
}

func NewDeleteBuildCmd(o *DeleteBuildOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: deleteAliases,
		Short:   "Delete build",
		Example: `
  # Delete build 'build1' in namespace 'ns1'
  knctl build delete -b build1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.BuildFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *DeleteBuildOptions) Run() error {
	buildClient, err := o.depsFactory.BuildClient()
	if err != nil {
		return err
	}

	err = buildClient.BuildV1alpha1().Builds(o.BuildFlags.NamespaceFlags.Name).Delete(o.BuildFlags.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Deleting build: %s", err)
	}

	// TODO idempotent?

	return nil
}
