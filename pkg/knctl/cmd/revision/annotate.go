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
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/spf13/cobra"
)

type AnnotateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	RevisionFlags cmdflags.RevisionFlags
	AnnotateFlags cmdflags.AnnotateFlags
}

func NewAnnotateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *AnnotateOptions {
	return &AnnotateOptions{ui: ui, depsFactory: depsFactory}
}

func NewAnnotateCmd(o *AnnotateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "Annotate revision",
		Example: `
  # Annotate revision 'rev1' in namespace 'ns1' with key and value
  knctl revision annotate -r rev1 -a key=value -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RevisionFlags.Set(cmd, flagsFactory)
	o.AnnotateFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *AnnotateOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewReference(o.RevisionFlags, tags, servingClient).Revision()
	if err != nil {
		return err
	}

	anns := ctlservice.NewAnnotations(servingClient)

	annotations, err := o.AnnotateFlags.AsMap()
	if err != nil {
		return err
	}

	return anns.Annotate(revision, annotations)
}
