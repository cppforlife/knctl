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
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

type AnnotateRevisionOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	RevisionFlags RevisionFlags
	AnnotateFlags AnnotateFlags
}

func NewAnnotateRevisionOptions(ui ui.UI, depsFactory DepsFactory) *AnnotateRevisionOptions {
	return &AnnotateRevisionOptions{ui: ui, depsFactory: depsFactory}
}

func NewAnnotateRevisionCmd(o *AnnotateRevisionOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revision",
		Aliases: revisionAliases,
		Short:   "Annotate revision",
		Example: `
  # Annotate revision 'rev1' in namespace 'ns1' with key and value
  knctl annotate revision -r rev1 -a key=value -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.RevisionFlags.Set(cmd, flagsFactory)
	o.AnnotateFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *AnnotateRevisionOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	tags := ctlservice.NewTags(servingClient)

	revision, err := NewRevisionReference(o.RevisionFlags, tags, servingClient).Revision()
	if err != nil {
		return err
	}

	patchJSON, err := o.annotationsAdditionPatchJSON()
	if err != nil {
		return err
	}

	_, err = servingClient.ServingV1alpha1().Revisions(revision.Namespace).Patch(revision.Name, types.MergePatchType, patchJSON)
	if err != nil {
		return fmt.Errorf("Annotating revision: %s", err)
	}

	return nil
}

func (o *AnnotateRevisionOptions) annotationsAdditionPatchJSON() ([]byte, error) {
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
