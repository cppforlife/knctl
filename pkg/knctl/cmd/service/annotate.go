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

package service

import (
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlkube "github.com/cppforlife/knctl/pkg/knctl/kube"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

type AnnotateOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ServiceFlags  cmdflags.ServiceFlags
	AnnotateFlags cmdflags.AnnotateFlags
}

func NewAnnotateOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *AnnotateOptions {
	return &AnnotateOptions{ui: ui, depsFactory: depsFactory}
}

func NewAnnotateCmd(o *AnnotateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "Annotate service",
		Example: `
  # Annotate service 'srv1' in namespace 'ns1' with key and value
  knctl service annotate -s srv1 -a key=value -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	o.AnnotateFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *AnnotateOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	anns := ctlkube.NewAnnotations(func(type_ types.PatchType, data []byte) error {
		_, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Patch(o.ServiceFlags.Name, type_, data)
		return err
	})

	annotations, err := o.AnnotateFlags.AsMap()
	if err != nil {
		return err
	}

	return anns.Add(annotations)
}
