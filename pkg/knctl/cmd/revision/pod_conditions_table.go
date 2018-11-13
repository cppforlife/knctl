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
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	"github.com/mitchellh/go-wordwrap"
	corev1 "k8s.io/api/core/v1"
)

type PodConditionsTable struct {
	podsCh chan corev1.Pod
}

func NewPodConditionsTable(podsCh chan corev1.Pod) PodConditionsTable {
	return PodConditionsTable{podsCh}
}

func (t PodConditionsTable) Print(ui ui.UI) {
	table := uitable.Table{
		Title: fmt.Sprintf("Pods conditions"),
		// TODO Content: "conditions",

		Header: []uitable.Header{
			uitable.NewHeader("Pod"),
			uitable.NewHeader("Type"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Age"),
			uitable.NewHeader("Reason"),
			uitable.NewHeader("Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for pod := range t.podsCh {
		for _, cond := range pod.Status.Conditions {
			table.Rows = append(table.Rows, []uitable.Value{
				uitable.NewValueString(pod.Name),
				uitable.NewValueString(string(cond.Type)),
				uitable.ValueFmt{
					V:     uitable.NewValueString(string(cond.Status)),
					Error: cond.Status != corev1.ConditionTrue,
				},
				cmdcore.NewValueAge(cond.LastTransitionTime.Time),
				uitable.NewValueString(cond.Reason),
				uitable.NewValueString(wordwrap.WrapString(cond.Message, 80)),
			})
		}
	}

	ui.PrintTable(table)
}
