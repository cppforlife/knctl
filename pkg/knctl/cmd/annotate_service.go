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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

type AnnotateServiceOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags  ServiceFlags
	AnnotateFlags AnnotateFlags
}

func NewAnnotateServiceOptions(ui ui.UI, depsFactory DepsFactory) *AnnotateServiceOptions {
	return &AnnotateServiceOptions{ui: ui, depsFactory: depsFactory}
}

func NewAnnotateServiceCmd(o *AnnotateServiceOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Aliases: serviceAliases,
		Short:   "Annotate service",
		Example: `
  # Annotate service 'srv1' in namespace 'ns1' with key and value
  knctl annotate service -s srv1 -a key=value -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd)
	o.AnnotateFlags.Set(cmd)
	return cmd
}

func (o *AnnotateServiceOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	patchJSON, err := o.annotationsAdditionPatchJSON()
	if err != nil {
		return err
	}

	_, err = servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Patch(o.ServiceFlags.Name, types.MergePatchType, patchJSON)
	if err != nil {
		return fmt.Errorf("Annotating service: %s", err)
	}

	return nil
}

func (o *AnnotateServiceOptions) annotationsAdditionPatchJSON() ([]byte, error) {
	annotations := map[string]interface{}{}

	for _, kv := range o.AnnotateFlags.Annotations {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected annotation to be in format 'KEY=VALUE'")
		}
		annotations[pieces[0]] = pieces[1]
	}

	mergePatch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": annotations,
		},
	}

	return json.Marshal(mergePatch)
}
