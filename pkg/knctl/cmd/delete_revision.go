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
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteRevisionOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	RevisionFlags RevisionFlags
}

func NewDeleteRevisionOptions(ui ui.UI, depsFactory DepsFactory) *DeleteRevisionOptions {
	return &DeleteRevisionOptions{ui: ui, depsFactory: depsFactory}
}

func NewDeleteRevisionCmd(o *DeleteRevisionOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revision",
		Short: "Delete revision",
		Example: `
  # Delete revision 'rev1' in namespace 'ns1'
  knctl delete revision -r rev1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RevisionFlags.Set(cmd)
	return cmd
}

func (o *DeleteRevisionOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewRevisionReference(o.RevisionFlags, tags, servingClient).Revision()
	if err != nil {
		return err
	}

	err = servingClient.ServingV1alpha1().Revisions(revision.Namespace).Delete(revision.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Deleting revision: %s", err)
	}

	// TODO idempotent?

	return nil
}
