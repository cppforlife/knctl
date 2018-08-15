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
	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
)

type DeployOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags ServiceFlags
	DeployFlags  DeployFlags
}

func NewDeployOptions(ui ui.UI, depsFactory DepsFactory) *DeployOptions {
	return &DeployOptions{ui: ui, depsFactory: depsFactory}
}

func NewDeployCmd(o *DeployOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy service",
		Example: `
  # Deploy service 'srv1' with a given image and one environment variable in namespace 'ns1'
  knctl deploy -s srv1 --image gcr.io/knative-samples/helloworld-go --env TARGET=123 -n ns1

  # Deploy service 'srv1' from Git repo and one environment variable in namespace 'ns1'
  knctl deploy -s srv1 --image gcr.io/your-account/your-image --git-url https://github.com/cppforlife/simple-app --git-revision master --env TARGET=123 -n ns1`, // TODO replace example
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd)
	o.DeployFlags.Set(cmd)
	return cmd
}

func (o *DeployOptions) Run() error {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return err
	}

	buildClient, err := o.depsFactory.BuildClient()
	if err != nil {
		return err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	restConfig, err := o.depsFactory.RESTConfig()
	if err != nil {
		return err
	}

	// TODO should we just automatically label it?
	err = NewIstio(coreClient).ExpectNamespaceToBeEnabled(o.ServiceFlags.NamespaceFlags.Name)
	if err != nil {
		return err
	}

	service, err := ServiceSpec{}.Build(o.ServiceFlags, o.DeployFlags)
	if err != nil {
		return err
	}

	buildObjFactory := ctlbuild.NewFactory(buildClient, coreClient, restConfig)
	serviceObj := ctlservice.NewService(service, servingClient, buildClient, coreClient, buildObjFactory)

	var lastRevision *v1alpha1.Revision

	// TODO find a better way to deal with generated service names
	if !o.DeployFlags.GenerateNameFlags.GenerateName {
		lastRevision, err = serviceObj.LastRevision()
		if err != nil {
			return err
		}
	}

	createdService, err := serviceObj.CreateOrUpdate()
	if err != nil {
		return err
	}

	o.printTable(createdService)

	tags := ctlservice.NewTags(servingClient)

	err = o.updateRevisionTags(serviceObj, tags, lastRevision)
	if err != nil {
		return err
	}

	if service.Spec.RunLatest.Configuration.Build != nil {
		cancelCh := make(chan struct{})

		buildObj, err := serviceObj.CreatedBuildSinceRevision(lastRevision)
		if err != nil {
			return err
		}

		if len(o.DeployFlags.BuildCreateArgsFlags.SourceDirectory) > 0 {
			err = buildObj.UploadSource(o.DeployFlags.BuildCreateArgsFlags.SourceDirectory, o.ui, cancelCh)
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

	return nil
}

func (o *DeployOptions) updateRevisionTags(
	serviceObj *ctlservice.Service, tags ctlservice.Tags, lastRevision *v1alpha1.Revision) error {

	if lastRevision != nil {
		o.ui.PrintLinef("Waiting for new revision (after revision '%s') to be created...", lastRevision.Name)
	} else {
		o.ui.PrintLinef("Waiting for new revision to be created...")
	}

	newLastRevision, err := serviceObj.CreatedRevisionSinceRevision(lastRevision)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Tagging new revision '%s' as '%s'", newLastRevision.Name, ctlservice.TagsLatest)

	err = tags.Repoint(newLastRevision, ctlservice.TagsLatest)
	if err != nil {
		return err
	}

	// If there was no last revision, let's tag new revision as previous
	if lastRevision != nil {
		o.ui.PrintLinef("Tagging older revision '%s' as '%s'", lastRevision.Name, ctlservice.TagsPrevious)

		err = tags.Repoint(lastRevision, ctlservice.TagsPrevious)
		if err != nil {
			return err
		}
	} else {
		o.ui.PrintLinef("Tagging new revision '%s' as '%s'", newLastRevision.Name, ctlservice.TagsPrevious)

		err = tags.Repoint(newLastRevision, ctlservice.TagsPrevious)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *DeployOptions) printTable(svc *v1alpha1.Service) {
	table := uitable.Table{
		Header: []uitable.Header{
			uitable.NewHeader("Name"),
		},

		Transpose: true,

		Rows: [][]uitable.Value{
			{uitable.NewValueString(svc.Name)},
		},
	}

	o.ui.PrintTable(table)
}
