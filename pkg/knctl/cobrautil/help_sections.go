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

package cobrautil

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cobra.AddTemplateFunc("commandsWithAnnotation", func(cmd *cobra.Command, key, value string) []*cobra.Command {
		var result []*cobra.Command
		for _, c := range cmd.Commands() {
			anns := map[string]string{}
			if c.Annotations != nil {
				anns = c.Annotations
			}
			if anns[key] == value {
				result = append(result, c)
			}
		}
		return result
	})
}

type HelpSection struct {
	Key   string
	Value string
	Title string
}

func HelpSectionsUsageTemplate(sections []HelpSection) string {
	unmodifiedCmd := &cobra.Command{}
	usageTemplate := unmodifiedCmd.UsageTemplate()

	const defaultTpl = `{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}`

	if !strings.Contains(usageTemplate, defaultTpl) {
		panic("Expected to find available commands section in spf13/cobra default usage template")
	}

	newTpl := "{{if .HasAvailableSubCommands}}"

	for _, section := range sections {
		newTpl += fmt.Sprintf(`{{$cmds := (commandsWithAnnotation . "%s" "%s")}}{{if $cmds}}

%s{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}`, section.Key, section.Value, section.Title)
	}

	newTpl += "{{end}}"

	return strings.Replace(usageTemplate, defaultTpl, newTpl, 1)
}
