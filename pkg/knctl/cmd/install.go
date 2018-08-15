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
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	istioURL        = "https://raw.githubusercontent.com/knative/serving/v0.1.0/third_party/istio-0.8.0/istio.yaml"
	knativeFullURL  = "https://github.com/knative/serving/releases/download/v0.1.0/release.yaml"
	knativeNoMonURL = "https://github.com/knative/serving/releases/download/v0.1.0/release-no-mon.yaml"
)

type InstallOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	ExcludeMonitoring bool

	kubeconfigFlags *KubeconfigFlags
}

func NewInstallOptions(ui ui.UI, depsFactory DepsFactory, kubeconfigFlags *KubeconfigFlags) *InstallOptions {
	return &InstallOptions{ui: ui, depsFactory: depsFactory, kubeconfigFlags: kubeconfigFlags}
}

func NewInstallCmd(o *InstallOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Knative and Istio",
		Long:  "Requires 'kubectl' command installed on a the system.",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	cmd.Flags().BoolVarP(&o.ExcludeMonitoring, "exclude-monitoring", "m", false, "Exclude installation of monitoring components")
	return cmd
}

func (o *InstallOptions) Run() error {
	// TODO remove kubectl dependency
	// TODO install latest dev version
	// TODO check kube versions is 1.10+
	// TODO grant cluster-admin permissions to the current user for Istio

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	istio := NewIstio(coreClient)
	knativeURL := knativeFullURL

	if o.ExcludeMonitoring {
		knativeURL = knativeNoMonURL
	}

	components := []InstallationComponent{
		{"Istio", istioURL, NamespaceReadiness{istio.SystemNamespaceName(), coreClient}, o.ui, o.kubeconfigFlags},
		{"Knative", knativeURL, NamespaceReadiness{"knative-serving", coreClient}, o.ui, o.kubeconfigFlags},
	}

	for _, c := range components {
		err = c.Install()
		if err != nil {
			return err
		}

		err = c.Monitor()
		if err != nil {
			return err
		}
	}

	o.ui.PrintLinef("Enabling Istio in namespace 'default'")

	return istio.EnableNamespace("default")
}

type InstallationComponent struct {
	Name string
	URL  string

	nsReadiness NamespaceReadiness
	ui          ui.UI

	kubeconfigFlags *KubeconfigFlags
}

func (c InstallationComponent) Install() error {
	c.ui.PrintLinef("Installing %s from '%s'", c.Name, c.URL)

	opts := []string{"--kubeconfig", c.kubeconfigFlags.Path, "apply", "-f", c.URL}

	_, err := exec.Command("kubectl", opts...).Output()
	if err != nil {
		return fmt.Errorf("Installing %s: %s", c.Name, err)
	}

	return nil
}

func (c InstallationComponent) Monitor() error {
	c.ui.PrintLinef("Waiting for %s to start...", c.Name)
	return c.nsReadiness.Monitor()
}

type NamespaceReadiness struct {
	Namespace  string
	coreClient kubernetes.Interface
}

func (n NamespaceReadiness) Monitor() error {
	var nonReadyPodNames []string

	for i := 0; i < 1000; i++ {
		allReady := true
		nonReadyPodNames = []string{}

		pods, err := n.coreClient.CoreV1().Pods(n.Namespace).List(metav1.ListOptions{})
		if err != nil {
			allReady = false
		}

		for _, pod := range pods.Items {
			if !(PodReadiness{pod}).IsRunningOrComplete() {
				allReady = false
				nonReadyPodNames = append(nonReadyPodNames, pod.Name)
			}
		}

		if allReady {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf(
		"Expected namespace '%s' to have all components be 'Running' or 'Completed' but found non-ready Pods: %s",
		n.Namespace, strings.Join(nonReadyPodNames, ", "))
}

type PodReadiness struct {
	Pod corev1.Pod
}

func (p PodReadiness) IsRunningOrComplete() bool {
	return p.Pod.Status.Phase == corev1.PodRunning || p.Pod.Status.Phase == corev1.PodSucceeded
}
