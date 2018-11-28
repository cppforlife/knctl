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

package build_test

import (
	"reflect"
	"testing"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	"github.com/knative/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildSpec(t *testing.T) {
	spec, err := ctlbuild.BuildSpec{}.Build(ctlbuild.BuildSpecOpts{
		ServiceAccountName: "test-service-account",
		GitURL:             "test-git-url",
		GitRevision:        "test-git-revision",
		Image:              "test-image",
	})
	if err != nil {
		t.Fatalf("Expected build spec to build successfully: %s", err)
	}

	expectedSpec := v1alpha1.BuildSpec{
		ServiceAccountName: "test-service-account",
		Source: &v1alpha1.SourceSpec{
			Git: &v1alpha1.GitSourceSpec{
				Url:      "test-git-url",
				Revision: "test-git-revision",
			},
		},
		Steps: []corev1.Container{
			{
				Name:  "build-and-push",
				Image: "gcr.io/kaniko-project/executor",
				Args: []string{
					"--dockerfile=/workspace/Dockerfile",
					"--destination=test-image",
				},
			},
		},
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expect spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}

func TestBuildSpecWithSourceDirectory(t *testing.T) {
	spec, err := ctlbuild.BuildSpec{}.Build(ctlbuild.BuildSpecOpts{
		ServiceAccountName: "test-service-account",
		SourceDirectory:    "test-source-dir",
		Image:              "test-image",
	})
	if err != nil {
		t.Fatalf("Expected build spec to build successfully: %s", err)
	}

	expectedSpec := v1alpha1.BuildSpec{
		ServiceAccountName: "test-service-account",
		Source: &v1alpha1.SourceSpec{
			Custom: &corev1.Container{
				Image:   "ubuntu:xenial",
				Command: []string{"/bin/bash"},
				Args: []string{
					"-c",
					"touch /tmp/SOURCE_UPLOAD_DONE; while [ -f /tmp/SOURCE_UPLOAD_DONE ]; do sleep 1; done",
				},
			},
		},
		Steps: []corev1.Container{
			{
				Name:  "build-and-push",
				Image: "gcr.io/kaniko-project/executor",
				Args: []string{
					"--dockerfile=/workspace/Dockerfile",
					"--destination=test-image",
				},
			},
		},
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expect spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}

func TestBuildSpecWithCustomBuildTemplate(t *testing.T) {
	spec, err := ctlbuild.BuildSpec{}.Build(ctlbuild.BuildSpecOpts{
		ServiceAccountName: "test-service-account",
		GitURL:             "test-git-url",
		GitRevision:        "test-git-revision",
		Template:           "test-template",
		TemplateArgs:       []string{"test-template-arg-key=test-template-arg-val"},
		TemplateEnv:        []string{"test-template-env-key=test-template-env-val"},
		Image:              "test-image",
	})
	if err != nil {
		t.Fatalf("Expected build spec to build successfully: %s", err)
	}

	expectedSpec := v1alpha1.BuildSpec{
		ServiceAccountName: "test-service-account",
		Source: &v1alpha1.SourceSpec{
			Git: &v1alpha1.GitSourceSpec{
				Url:      "test-git-url",
				Revision: "test-git-revision",
			},
		},
		Template: &v1alpha1.TemplateInstantiationSpec{
			Name: "test-template",
			Kind: "BuildTemplate",
			Arguments: []v1alpha1.ArgumentSpec{
				{
					Name:  "test-template-arg-key",
					Value: "test-template-arg-val",
				},
				{
					Name:  "IMAGE",
					Value: "test-image",
				},
			},
			Env: []corev1.EnvVar{{
				Name:  "test-template-env-key",
				Value: "test-template-env-val",
			}},
		},
	}

	if !reflect.DeepEqual(spec.Template, expectedSpec.Template) {
		t.Fatalf("Expect spec.template '%#v' to equal '%#v'", spec.Template, expectedSpec.Template)
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expect spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}
