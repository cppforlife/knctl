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
	gosha256 "crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	istioAsset = YAMLAsset{
		URL:    "https://raw.githubusercontent.com/knative/serving/v0.1.0/third_party/istio-0.8.0/istio.yaml",
		SHA256: "22b8ea622472d6253bf713d5767e9e179091e55e7e5aaf8fa0a586e5b57deb93",
	}
	knativeFullAsset = YAMLAsset{
		URL:    "https://github.com/knative/serving/releases/download/v0.1.0/release.yaml",
		SHA256: "6e92473ed675003fd52cde4ba9dc409073dbc5d7259bd496c88fb5b99633795e",
	}
	knativeNoMonAsset = YAMLAsset{
		URL:    "https://github.com/knative/serving/releases/download/v0.1.0/release-no-mon.yaml",
		SHA256: "69408b14133242e4daa8f62719eab998200e3b453555e12b678c4ffade43e332",
	}
)

type InstallOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	NodePorts         bool
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
		Long: `Install Knative and Istio.

Requires 'kubectl' command installed on a the system.`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	cmd.Flags().BoolVarP(&o.NodePorts, "node-ports", "p", false, "Use service type NodePorts instead of type LoadBalancer")
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
	knativeAsset := knativeFullAsset

	if o.ExcludeMonitoring {
		knativeAsset = knativeNoMonAsset
	}

	components := []InstallationComponent{
		{"Istio", YAMLSource{istioAsset, o.NodePorts}, NamespaceReadiness{istio.SystemNamespaceName(), coreClient}, o.ui, o.kubeconfigFlags},
		{"Knative", YAMLSource{knativeAsset, o.NodePorts}, NamespaceReadiness{"knative-serving", coreClient}, o.ui, o.kubeconfigFlags},
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
	Name   string
	source YAMLSource

	nsReadiness NamespaceReadiness
	ui          ui.UI

	kubeconfigFlags *KubeconfigFlags
}

func (c InstallationComponent) Install() error {
	c.ui.PrintLinef("Installing %s from '%s'", c.Name, c.source.Source())

	content, err := c.source.Content()
	if err != nil {
		return err
	}

	opts := []string{"--kubeconfig", c.kubeconfigFlags.Path, "apply", "-f", "-"}

	cmd := exec.Command("kubectl", opts...)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = uiLinesWriter{c.ui}
	cmd.Stderr = uiLinesWriter{c.ui}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Installing %s: %s", c.Name, err)
	}

	return nil
}

func (c InstallationComponent) Monitor() error {
	c.ui.PrintLinef("Waiting for %s to start...", c.Name)
	return c.nsReadiness.Monitor()
}

type YAMLSource struct {
	Asset     YAMLAsset
	NodePorts bool
}

func (s YAMLSource) Source() string {
	return s.Asset.URL
}

func (s YAMLSource) Content() (string, error) {
	content, err := s.download(s.Asset.URL, s.Asset.SHA256)
	if err != nil {
		return "", err
	}

	if s.NodePorts {
		content = strings.Replace(content, "type: LoadBalancer", "type: NodePort", -1)
	}

	return content, nil
}

func (YAMLSource) download(url, expectedSHA256 string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Downloading YAML from URL '%s': %s", url, err)
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Reading YAML from URL '%s': %s", url, err)
	}

	if fmt.Sprintf("%x", gosha256.Sum256(result)) != expectedSHA256 {
		return "", fmt.Errorf("Expected URL '%s' content to match SHA256 '%s' but did not", url, expectedSHA256)
	}

	return string(result), nil
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

type YAMLAsset struct {
	URL    string
	SHA256 string
}

type uiLinesWriter struct {
	ui ui.UI
}

var _ io.Writer = uiLinesWriter{}

func (w uiLinesWriter) Write(p []byte) (n int, err error) {
	w.ui.BeginLinef("%s", p)
	return len(p), nil
}
