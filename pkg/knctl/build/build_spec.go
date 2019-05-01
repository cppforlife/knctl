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
	"strings"
	"time"

	"github.com/knative/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	buildSpecImageArgName = "IMAGE"
)

type BuildSpec struct{}

type BuildSpecOpts struct {
	SourceDirectory string

	GitURL      string
	GitRevision string

	ServiceAccountName string

	TemplateName string
	TemplateKind string
	TemplateArgs []string
	TemplateEnv  []string

	BuildArgs []string

	Image   string
	Timeout time.Duration
}

func (s BuildSpec) Build(opts BuildSpecOpts) (v1alpha1.BuildSpec, error) {
	source, err := s.source(opts)
	if err != nil {
		return v1alpha1.BuildSpec{}, err
	}

	template, err := s.template(opts)
	if err != nil {
		return v1alpha1.BuildSpec{}, err
	}

	steps := s.nonTemplateSteps(opts)

	spec := v1alpha1.BuildSpec{
		ServiceAccountName: opts.ServiceAccountName,
		Source:             source,
		Template:           template,
		Steps:              steps,
	}

	if opts.Timeout.Nanoseconds() > 0 {
		spec.Timeout = &metav1.Duration{opts.Timeout}
	}

	return spec, nil
}

func (s BuildSpec) source(opts BuildSpecOpts) (*v1alpha1.SourceSpec, error) {
	switch {
	case len(opts.SourceDirectory) > 0:
		return &v1alpha1.SourceSpec{
			Custom: &corev1.Container{
				Image:   "ubuntu:xenial",
				Command: []string{"/bin/bash"},
				Args: []string{
					"-c",
					fmt.Sprintf("touch %s; while [ -f %s ]; do sleep 1; done",
						clusterBuilderCustomSourceTriggerFile, clusterBuilderCustomSourceTriggerFile),
				},
			},
		}, nil

	case len(opts.GitURL) > 0:
		return &v1alpha1.SourceSpec{
			Git: &v1alpha1.GitSourceSpec{
				Url:      opts.GitURL,
				Revision: opts.GitRevision,
			},
		}, nil

	default:
		return nil, fmt.Errorf("Unknown build source")
	}
}

func (s BuildSpec) template(opts BuildSpecOpts) (*v1alpha1.TemplateInstantiationSpec, error) {
	if len(opts.TemplateName) == 0 {
		return nil, nil
	}

	args, err := s.templateArgs(opts)
	if err != nil {
		return nil, err
	}

	env, err := s.templateEnv(opts)
	if err != nil {
		return nil, err
	}

	templateKind := v1alpha1.BuildTemplateKind
	if opts.TemplateKind == "cluster" || opts.TemplateKind == "Cluster" || opts.TemplateKind == "ClusterBuildTemplate" {
		templateKind = v1alpha1.ClusterBuildTemplateKind
	}

	return &v1alpha1.TemplateInstantiationSpec{
		Name:      opts.TemplateName,
		Kind:      templateKind,
		Arguments: args,
		Env:       env,
	}, nil
}

func (s BuildSpec) nonTemplateSteps(opts BuildSpecOpts) []corev1.Container {
	if len(opts.TemplateName) > 0 {
		return nil
	}

	args := []string{
		"--dockerfile=/workspace/Dockerfile",
		"--destination=" + opts.Image,
	}

	for _, kv := range opts.BuildArgs {
		args = append(args, "--build-arg="+kv)
	}

	return []corev1.Container{
		{
			Name:  "build-and-push",
			Image: "gcr.io/kaniko-project/executor",
			Args:  args,
		},
	}
}

func (s BuildSpec) templateArgs(opts BuildSpecOpts) ([]v1alpha1.ArgumentSpec, error) {
	var args []v1alpha1.ArgumentSpec
	var hadImageArg bool

	for _, kv := range opts.TemplateArgs {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected build template argument to be in format 'key=value'")
		}

		args = append(args, v1alpha1.ArgumentSpec{
			Name:  pieces[0],
			Value: pieces[1],
		})

		if pieces[0] == buildSpecImageArgName {
			hadImageArg = true
		}
	}

	if !hadImageArg {
		args = append(args, v1alpha1.ArgumentSpec{
			Name:  buildSpecImageArgName,
			Value: opts.Image,
		})
	}

	return args, nil
}

func (s BuildSpec) templateEnv(opts BuildSpecOpts) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	for _, kv := range opts.TemplateEnv {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected build template environment variable to be in format 'key=value'")
		}

		env = append(env, corev1.EnvVar{
			Name:  pieces[0],
			Value: pieces[1],
		})
	}

	return env, nil
}
