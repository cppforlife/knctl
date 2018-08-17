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

type UIFlags struct {
	TTY            bool
	NoColor        bool
	JSON           bool
	NonInteractive bool
	Columns        []string
}

func (f *UIFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	cmd.PersistentFlags().BoolVar(&f.TTY, "tty", false, "Force TTY-like output")
	cmd.PersistentFlags().BoolVar(&f.NoColor, "no-color", false, "Disable colorized output")
	cmd.PersistentFlags().BoolVar(&f.JSON, "json", false, "Output as JSON")
	cmd.PersistentFlags().BoolVar(&f.NonInteractive, "non-interactive", false, "Don't ask for user input")
	cmd.PersistentFlags().StringSliceVar(&f.Columns, "column", nil, "Filter to show only given columns")
}

func (f *UIFlags) ConfigureUI(ui *ui.ConfUI) {
	ui.EnableTTY(f.TTY)

	if !f.NoColor {
		ui.EnableColor()
	}

	if f.JSON {
		ui.EnableJSON()
	}

	if f.NonInteractive {
		ui.EnableNonInteractive()
	}

	if len(f.Columns) > 0 {
		headers := []uitable.Header{}
		for _, col := range f.Columns {
			headers = append(headers, uitable.Header{
				Key:    uitable.KeyifyHeader(col),
				Hidden: false,
			})
		}

		ui.ShowColumns(headers)
	}
}
