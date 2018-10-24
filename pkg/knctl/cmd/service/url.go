/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Open 2.0 (the "License");
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
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type URLOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ServiceFlags cmdflags.ServiceFlags
	CurlFlags    CurlFlags
}

func NewURLOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *URLOptions {
	return &URLOptions{ui: ui, depsFactory: depsFactory}
}

func NewURLCmd(o *URLOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url",
		Short: "Print service URL",
		Example: `
  # Print service 'svc1' URL in namespace 'ns1'
  knctl service url -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	o.CurlFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *URLOptions) Run() error {
	url, err := o.addr()
	if err != nil {
		return err
	}

	o.ui.PrintBlock([]byte(url))

	return nil
}

func (o *URLOptions) addr() (string, error) {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return "", err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return "", err
	}

	url, err := ServiceAddress{service, coreClient}.URL(o.CurlFlags.Port, true)
	if err != nil {
		return "", err
	}

	return url, nil
}
