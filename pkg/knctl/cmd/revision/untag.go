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

type UntagOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	RevisionFlags cmdflags.RevisionFlags
	TagFlags      TagFlags
}

func NewUntagOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *UntagOptions {
	return &UntagOptions{ui: ui, depsFactory: depsFactory}
}

func NewUntagCmd(o *UntagOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "untag",
		Short: "Untag revision",
		Example: `
  # Untag revision 'rev1' in namespace 'ns1' as 'stable'
  knctl revision untag -r rev1 -t stable -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RevisionFlags.Set(cmd, flagsFactory)
	o.TagFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *UntagOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewReference(o.RevisionFlags, tags, servingClient).Revision()
	if err != nil {
		return err
	}

	for _, tag := range o.TagFlags.Tags {
		err := tags.Untag(*revision, tag)
		if err != nil {
			return err
		}
	}

	return nil
}
