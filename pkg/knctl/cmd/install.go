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
	"strconv"
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	InstallIstioAsset = InstallationAsset{
		URL:    "https://raw.githubusercontent.com/knative/serving/38c0d500fcd4a65b24b103b54bedf6dacc985170/third_party/istio-1.0.2/istio.yaml",
		SHA256: "92377c1600653bddb7e0e0a3a481e15d4f193f3dee36bb30399cc9aba7d628bf",
	}
	InstallKnativeFullAsset = InstallationAsset{
		URL:    "https://github.com/knative/serving/releases/download/v0.1.1/release.yaml",
		SHA256: "81d619b995ee36650ac4fe5ba54705cde569a92457aee18a03a8a45e5a9b8b77",
	}
	InstallKnativeNoMonAsset = InstallationAsset{
		URL:    "https://github.com/knative/serving/releases/download/v0.1.1/release-no-mon.yaml",
		SHA256: "db82bf221513bf5738bec694f3654df6111d74ad5cd1f3d69cb25422755437a7",
	}
)

type InstallOptions struct {
	ui          ui.UI
	depsFactory DepsFactory

	NodePorts         bool
	ExcludeMonitoring bool
	VersionCheck      bool

	kubeconfigFlags *KubeconfigFlags
}

func NewInstallOptions(ui ui.UI, depsFactory DepsFactory, kubeconfigFlags *KubeconfigFlags) *InstallOptions {
	return &InstallOptions{ui: ui, depsFactory: depsFactory, kubeconfigFlags: kubeconfigFlags}
}

func NewInstallCmd(o *InstallOptions, flagsFactory FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Knative and Istio",
		Long: `Install Knative and Istio.

Requires 'kubectl' command installed on a the system.`,
		Annotations: map[string]string{
			systemGroup.Key: systemGroup.Value,
		},
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	cmd.Flags().BoolVarP(&o.NodePorts, "node-ports", "p", false, "Use service type NodePorts instead of type LoadBalancer")
	cmd.Flags().BoolVarP(&o.ExcludeMonitoring, "exclude-monitoring", "m", false, "Exclude installation of monitoring components")
	cmd.Flags().BoolVar(&o.VersionCheck, "version-check", true, "Check minimum Kubernetes API server version")
	return cmd
}

func (o *InstallOptions) Run() error {
	// TODO remove kubectl dependency
	// TODO install latest dev version
	// TODO grant cluster-admin permissions to the current user for Istio

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	if o.VersionCheck {
		err = o.ensureMinimumServerVersion(coreClient)
		if err != nil {
			return fmt.Errorf("%s (skip via --version-check=false)", err)
		}
	}

	istio := NewIstio()
	knativeAsset := InstallKnativeFullAsset

	if o.ExcludeMonitoring {
		knativeAsset = InstallKnativeNoMonAsset
	}

	components := []InstallationComponent{
		{"Istio", YAMLSource{InstallIstioAsset, o.NodePorts}, NamespaceReadiness{istio.SystemNamespaceName(), o.ui, coreClient}, o.ui, o.kubeconfigFlags, 1},
		{"Knative", YAMLSource{knativeAsset, o.NodePorts}, NamespaceReadiness{"knative-serving", o.ui, coreClient}, o.ui, o.kubeconfigFlags, 0},
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

	return nil
}

func (o *InstallOptions) ensureMinimumServerVersion(coreClient kubernetes.Interface) error {
	version, err := coreClient.Discovery().ServerVersion()
	if err != nil {
		return err
	}

	majorI, err := strconv.Atoi(version.Major)
	if err != nil {
		return fmt.Errorf("Converting major version '%s' to int: %s", version.Major, err)
	}

	// GKE shows minor as "10+"
	minorI, err := strconv.Atoi(strings.TrimRight(version.Minor, "-+"))
	if err != nil {
		return fmt.Errorf("Converting minor version '%s' to int: %s", version.Minor, err)
	}

	if majorI == 1 && minorI < 10 {
		return fmt.Errorf("Expected Kubernetes API server version to be >=1.10")
	}

	return nil
}

type InstallationComponent struct {
	Name   string
	source YAMLSource

	nsReadiness NamespaceReadiness
	ui          ui.UI

	kubeconfigFlags *KubeconfigFlags
	retryCount      int
}

func (c InstallationComponent) Install() error {
	c.ui.PrintLinef("Installing %s from '%s'", c.Name, c.source.Source())

	kubeconfigPath, err := c.kubeconfigFlags.Path.Value()
	if err != nil {
		return err
	}

	content, err := c.source.Content()
	if err != nil {
		return err
	}

	var lastErr error

	for i := 0; i < c.retryCount+1; i++ {
		lastErr = c.runKubectl(content, []string{"--kubeconfig", kubeconfigPath, "apply", "-f", "-"})
		if lastErr == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("Installing %s: %s", c.Name, lastErr)
}

func (c InstallationComponent) runKubectl(content string, opts []string) error {
	cmd := exec.Command("kubectl", opts...)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = uiLinesWriter{c.ui}
	cmd.Stderr = uiLinesWriter{c.ui}
	return cmd.Run()
}

func (c InstallationComponent) Monitor() error {
	c.ui.PrintLinef("Waiting for %s to start...", c.Name)
	return c.nsReadiness.Monitor()
}

type YAMLSource struct {
	Asset     InstallationAsset
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
	ui         ui.UI
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

		if i%20 == 0 {
			n.ui.PrintLinef("Non-ready pods: %s", strings.Join(nonReadyPodNames, ", "))
		}
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

type InstallationAsset struct {
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
