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

package knative

import (
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	ctling "github.com/cppforlife/knctl/pkg/knctl/ingress"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type UninstallOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	ExcludeMonitoring bool
}

func NewUninstallOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *UninstallOptions {
	return &UninstallOptions{ui: ui, depsFactory: depsFactory}
}

func NewUninstallCmd(o *UninstallOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Knative and Istio",
		Long: `Uninstall Knative and Istio. 

Requires 'kubectl' command installed on a the system.`,
		Annotations: map[string]string{
			cmdcore.SystemHelpGroup.Key: cmdcore.SystemHelpGroup.Value,
		},
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	return cmd
}

func (o *UninstallOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	istio := ctling.NewIstio()

	components := []UninstallationComponent{
		{"Knative Build", NamespaceRemoval{"knative-build", coreClient}, o.ui},
		{"Knative Serving", NamespaceRemoval{"knative-serving", coreClient}, o.ui},
		{"Knative Monitoring", NamespaceRemoval{"knative-monitoring", coreClient}, o.ui},
		{"Istio", NamespaceRemoval{istio.SystemNamespaceName(), coreClient}, o.ui},
	}

	for _, c := range components {
		err = c.Uninstall()
		if err != nil {
			return err
		}
	}

	return nil
}

type UninstallationComponent struct {
	Name string

	nsRemoval NamespaceRemoval
	ui        ui.UI
}

func (c UninstallationComponent) Uninstall() error {
	c.ui.PrintLinef("Uninstalling %s", c.Name)

	err := c.nsRemoval.Remove()
	if err != nil {
		return err
	}

	return c.Monitor()
}

func (c UninstallationComponent) Monitor() error {
	c.ui.PrintLinef("Waiting for namespace '%s' to be deleted...", c.nsRemoval.Namespace)
	return c.nsRemoval.Monitor()
}

type NamespaceRemoval struct {
	Namespace  string
	coreClient kubernetes.Interface
}

func (n NamespaceRemoval) Remove() error {
	err := n.coreClient.CoreV1().Namespaces().Delete(n.Namespace, &metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (n NamespaceRemoval) Monitor() error {
	for i := 0; i < 1000; i++ {
		namespaces, err := n.coreClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			return err
		}

		var found bool
		for _, namespace := range namespaces.Items {
			if namespace.Name == n.Namespace {
				found = true
				break
			}
		}

		if !found {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("Expected namespace '%s' to be deleted", n.Namespace)
}
