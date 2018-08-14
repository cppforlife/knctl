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

package build

import (
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	ctlkube "github.com/cppforlife/knctl/pkg/knctl/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	clusterBuilderCustomSourceStep        = "build-step-custom-source"
	clusterBuilderCustomSourceTriggerFile = "/tmp/SOURCE_UPLOAD_DONE"
)

type ClusterBuilderSource struct {
	dirPath    string
	waiter     BuildWaiter
	coreClient kubernetes.Interface
	restConfig *rest.Config
}

var _ Source = ClusterBuilderSource{}

func NewClusterBuilderSource(
	dirPath string,
	waiter BuildWaiter,
	coreClient kubernetes.Interface,
	restConfig *rest.Config,
) ClusterBuilderSource {
	return ClusterBuilderSource{dirPath, waiter, coreClient, restConfig}
}

func (s ClusterBuilderSource) Upload(ui ui.UI, cancelCh chan struct{}) error { // TODO cancel
	ui.PrintLinef("[%s] Uploading source code...", time.Now().Format(time.RFC3339))

	defer func() {
		ui.PrintLinef("[%s] Finished uploading source code...", time.Now().Format(time.RFC3339))
	}()

	build, err := s.waiter.WaitForClusterBuilderPodAssignment(cancelCh)
	if err != nil {
		return fmt.Errorf("Waiting for build to be assigned a pod: %s", err)
	}

	if build.Status.Cluster == nil {
		return fmt.Errorf("Expected build to have cluster configuration assigned")
	}

	podsClient := s.coreClient.CoreV1().Pods(build.Status.Cluster.Namespace)

	pod, err := podsClient.Get(build.Status.Cluster.PodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Getting assigned building pod: %s", err)
	}

	err = PodInitContainerRunningWatcher{*pod, podsClient, clusterBuilderCustomSourceStep}.Wait(cancelCh)
	if err != nil {
		return fmt.Errorf("Waiting for init container: %s", err)
	}

	executor := ctlkube.NewExec(*pod, clusterBuilderCustomSourceStep, s.coreClient, s.restConfig)

	err = ctlkube.NewDirCp(executor).Execute(s.dirPath, "/workspace")
	if err != nil {
		return fmt.Errorf("Uploading files: %s", err)
	}

	// TODO is there a race with touch?
	execArgs := []string{
		"/bin/bash", "-c", fmt.Sprintf("rm %s", clusterBuilderCustomSourceTriggerFile),
	}

	err = executor.Execute(execArgs, nil)
	if err != nil {
		return fmt.Errorf("Finishing upload: %s", err)
	}

	return nil
}
