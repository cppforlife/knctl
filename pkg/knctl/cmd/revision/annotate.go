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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/cppforlife/knctl/pkg/knctl/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
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

	patchJSON, err := o.annotationsAdditionPatchJSON()
	if err != nil {
		return err
	}

	return util.Retry(time.Second, 10*time.Second, func() (bool, error) {
		_, err := servingClient.ServingV1alpha1().Revisions(revision.Namespace).Patch(revision.Name, types.MergePatchType, patchJSON)
		if err != nil {
			return false, fmt.Errorf("Annotating revision: %s", err)
		}

		return true, nil
	})
}

func (o *AnnotateOptions) annotationsAdditionPatchJSON() ([]byte, error) {
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
