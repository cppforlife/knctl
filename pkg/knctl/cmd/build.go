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
	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildOptions struct {
	ui            ui.UI
	configFactory ConfigFactory
	depsFactory   DepsFactory
	cancelSignals CancelSignals

	BuildFlags       BuildFlags
	BuildCreateFlags BuildCreateFlags
}

func NewBuildOptions(
	ui ui.UI,
	configFactory ConfigFactory,
	depsFactory DepsFactory,
	cancelSignals CancelSignals,
) *BuildOptions {
	return &BuildOptions{ui: ui, configFactory: configFactory, depsFactory: depsFactory, cancelSignals: cancelSignals}
}

func NewBuildCmd(o *BuildOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Build source code into image",
		Example: `
  # Build Git repository into an image in namespace 'ns1'
  knctl build create -b build1 --git-url github.com/cppforlife/simple-app --git-revision master -i docker.io/cppforlife/simple-app -n ns1

  # Build from local source code in namespace 'ns1'
  # ( related: https://github.com/cppforlife/knctl/blob/master/docs/deploy-source-directory.md )
  knctl build create -b build1 -d=. -i index.docker.io/your-account/your-image --service-account serv-acct1 -n ns1

  # Build with custom build template in namespace 'ns1'
  # ( related: https://github.com/cppforlife/knctl/blob/master/docs/deploy-custom-build-template.md )
  knctl build create -b build1 -n ns1 \
      --git-url https://github.com/cppforlife/simple-app --git-revision master \
      --template buildpack --template-env GOPACKAGENAME=main \
      --service-account serv-acct1 --image index.docker.io/your-account/your-image`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.BuildFlags.Set(cmd, flagsFactory)
	o.BuildCreateFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *BuildOptions) Run() error {
	buildClient, err := o.depsFactory.BuildClient()
	if err != nil {
		return err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	restConfig, err := o.configFactory.RESTConfig()
	if err != nil {
		return err
	}

	build := &v1alpha1.Build{
		ObjectMeta: o.BuildCreateFlags.GenerateNameFlags.Apply(metav1.ObjectMeta{
			Name:      o.BuildFlags.Name,
			Namespace: o.BuildFlags.NamespaceFlags.Name,
		}),
	}

	build.Spec, err = ctlbuild.BuildSpec{}.Build(o.BuildCreateFlags.BuildSpecOpts)
	if err != nil {
		return err
	}

	createdBuild, err := buildClient.BuildV1alpha1().Builds(o.BuildFlags.NamespaceFlags.Name).Create(build)
	if err != nil {
		return err // TODO allow updating build?
	}

	o.printTable(createdBuild)

	cancelCh := make(chan struct{})
	o.cancelSignals.Watch(func() { close(cancelCh) })

	buildObjFactory := ctlbuild.NewFactory(buildClient, coreClient, restConfig)
	buildObj := buildObjFactory.New(createdBuild)

	if len(o.BuildCreateFlags.BuildCreateArgsFlags.SourceDirectory) > 0 {
		err = buildObj.UploadSource(o.BuildCreateFlags.BuildCreateArgsFlags.SourceDirectory, o.ui, cancelCh)
		if err != nil {
			return err
		}
	}

	err = buildObj.TailLogs(o.ui, cancelCh)
	if err != nil {
		return err
	}

	return buildObj.Error(cancelCh)
}

func (o *BuildOptions) printTable(b *v1alpha1.Build) {
	table := uitable.Table{
		Header: []uitable.Header{
			uitable.NewHeader("Name"),
		},

		Transpose: true,

		Rows: [][]uitable.Value{
			{uitable.NewValueString(b.Name)},
		},
	}

	o.ui.PrintTable(table)
}
