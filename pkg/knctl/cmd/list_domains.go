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
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
)

type ListDomainsOptions struct {
	ui          ui.UI
	depsFactory DepsFactory
}

func NewListDomainsOptions(ui ui.UI, depsFactory DepsFactory) *ListDomainsOptions {
	return &ListDomainsOptions{ui: ui, depsFactory: depsFactory}
}

func NewListDomainsCmd(o *ListDomainsOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domains",
		Aliases: []string{"dom", "domain"},
		Short:   "List domains",
		Long:    "List all domains",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	return cmd
}

func (o *ListDomainsOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	domains, err := NewDomains(coreClient).List()
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: "Domains",

		Content: "domains",

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Default"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
		},
	}

	for _, domain := range domains {
		table.Rows = append(table.Rows, []uitable.Value{
			uitable.NewValueString(domain.Name),
			uitable.NewValueBool(domain.Default),
		})
	}

	o.ui.PrintTable(table)

	return nil
}
