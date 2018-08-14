/*
Copyright 2016 The Kubernetes Authors.

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

package kube

import (
	"bytes"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Exec struct {
	pod        corev1.Pod
	container  string
	coreClient kubernetes.Interface
	restConfig *rest.Config
}

func NewExec(pod corev1.Pod, container string, coreClient kubernetes.Interface, restConfig *rest.Config) Exec {
	return Exec{pod, container, coreClient, restConfig}
}

func (s Exec) Execute(cmd []string, stdin io.Reader) error {
	req := s.coreClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(s.pod.Name).
		Namespace(s.pod.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Stderr:    true,
		Stdin:     stdin != nil,
		TTY:       false,
		Command:   cmd,
		Container: s.container,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(s.restConfig, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("Building executor: %s", err)
	}

	var stderr bytes.Buffer

	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return fmt.Errorf("Execution error: %s (stderr: %s)", err, stderr.String())
	}

	return nil
}
