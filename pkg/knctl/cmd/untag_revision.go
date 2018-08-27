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
	"github.com/cppforlife/go-cli-ui/ui"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/spf13/cobra"
)

type UntagRevisionOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	RevisionFlags RevisionFlags
	TagFlags      TagFlags
}

func NewUntagRevisionOptions(ui ui.UI, depsFactory DepsFactory) *UntagRevisionOptions {
	return &UntagRevisionOptions{ui: ui, depsFactory: depsFactory}
}

func NewUntagRevisionCmd(o *UntagRevisionOptions, flagsFactory FlagsFactory) *cobra.Command {
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

func (o *UntagRevisionOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewRevisionReference(o.RevisionFlags, tags, servingClient).Revision()
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
