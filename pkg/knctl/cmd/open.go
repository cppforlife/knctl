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

package cmd

import (
	"fmt"
	"os/exec"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OpenOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ServiceFlags ServiceFlags
}

func NewOpenOptions(ui ui.UI, depsFactory DepsFactory) *OpenOptions {
	return &OpenOptions{ui: ui, depsFactory: depsFactory}
}

func NewOpenCmd(o *OpenOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open web browser pointing at a service domain",
		Long: `Open web browser pointing at a service domain.

Requires 'open' command installed on the system.`,
		Example: `
# Open web browser pointing at service 'svc1' in namespace 'ns1'
knctl open -s svc1 -n ns1`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.ServiceFlags.Set(cmd, flagsFactory)
	return cmd
}

func (o *OpenOptions) Run() error {
	serviceDomain, err := o.serviceDomain()
	if err != nil {
		return err
	}

	// TODO Determine protocol for the entrypoint
	cmd := exec.Command("open", []string{"https://" + serviceDomain}...)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Starting browser: %s", err)
	}

	return nil
}

func (o *OpenOptions) serviceDomain() (string, error) {
	servingClient, err := o.depsFactory.ServingClient()
	if err != nil {
		return "", err
	}

	service, err := servingClient.ServingV1alpha1().Services(o.ServiceFlags.NamespaceFlags.Name).Get(o.ServiceFlags.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(service.Status.Domain) == 0 {
		return "", fmt.Errorf("Expected service '%s' to have non-empty domain", o.ServiceFlags.Name)
	}

	return service.Status.Domain, nil
}
