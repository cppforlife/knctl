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

package service_test

import (
	"reflect"
	"testing"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	cmdbld "github.com/cppforlife/knctl/pkg/knctl/cmd/build"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	. "github.com/cppforlife/knctl/pkg/knctl/cmd/service"
	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceSpecWithBuildConfiguration(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		BuildCreateArgsFlags: cmdbld.CreateArgsFlags{
			ctlbuild.BuildSpecOpts{
				GitURL:             "test-git-url",
				GitRevision:        "test-git-revision",
				ServiceAccountName: "test-service-account",
			},
		},
		Image:   "test-image",
		EnvVars: []string{"test-env-key1=test-env-val1"},

		RemoveKnctlDeployEnvVar: true,
	}

	spec, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err != nil {
		t.Fatalf("Expected error to not happen: %s", err)
	}

	expectedSpec := v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					Build: &v1alpha1.RawExtension{
						BuildSpec: &buildv1alpha1.BuildSpec{
							ServiceAccountName: "test-service-account",
							Source: &buildv1alpha1.SourceSpec{
								Git: &buildv1alpha1.GitSourceSpec{
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
						},
					},
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							ServiceAccountName: "test-service-account",
							Container: corev1.Container{
								Image: deployFlags.Image,
								Env:   []corev1.EnvVar{{Name: "test-env-key1", Value: "test-env-val1"}},
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expected spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}

func TestServiceSpecWithoutBuildConfiguration(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		Image:   "test-image",
		EnvVars: []string{"test-env-key1=test-env-val1"},

		RemoveKnctlDeployEnvVar: true,
	}

	spec, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err != nil {
		t.Fatalf("Expected error to not happen: %s", err)
	}

	expectedSpec := v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					Build: &v1alpha1.RawExtension{},
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: deployFlags.Image,
								Env:   []corev1.EnvVar{{Name: "test-env-key1", Value: "test-env-val1"}},
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expected spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}

func TestServiceSpecWithInvalidEnv(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		Image:   "test-image",
		EnvVars: []string{"test-env-key1"},

		RemoveKnctlDeployEnvVar: true,
	}

	_, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err == nil {
		t.Fatalf("Expected error to happen")
	}

	if err.Error() != "Expected environment variable to be in format 'ENV_KEY=value'" {
		t.Fatalf("Expected error to happen, but was '%s'", err)
	}
}

func TestServiceSpecWithMultipleEnv(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		Image:         "test-image",
		EnvVars:       []string{"test-env-key1=test-env-val1", "test-env-key2=test-env-val2"},
		EnvSecrets:    []string{"test-env-key3=test-secret1/key", "test-env-key4=test-secret2/key"},
		EnvConfigMaps: []string{"test-env-key5=test-config-map1/key", "test-env-key6=test-config-map2/key"},

		RemoveKnctlDeployEnvVar: true,
	}

	spec, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err != nil {
		t.Fatalf("Expected error to not happen: %s", err)
	}

	expectedSpec := v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					Build: &v1alpha1.RawExtension{},
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: deployFlags.Image,
								Env: []corev1.EnvVar{
									{Name: "test-env-key1", Value: "test-env-val1"},
									{Name: "test-env-key2", Value: "test-env-val2"},
									{Name: "test-env-key3", ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "test-secret1",
											},
											Key: "key",
										},
									}},
									{Name: "test-env-key4", ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "test-secret2",
											},
											Key: "key",
										},
									}},
									{Name: "test-env-key5", ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "test-config-map1",
											},
											Key: "key",
										},
									}},
									{Name: "test-env-key6", ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "test-config-map2",
											},
											Key: "key",
										},
									}},
								},
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(spec, expectedSpec) {
		t.Fatalf("Expected spec '%#v' to equal '%#v'", spec, expectedSpec)
	}
}

func TestServiceSpecWithInvalidEnvSecret(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		Image:      "test-image",
		EnvSecrets: []string{"test-env-secret-key1"},

		RemoveKnctlDeployEnvVar: true,
	}

	_, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err == nil {
		t.Fatalf("Expected error to happen")
	}

	if err.Error() != "Expected environment variable from secret to be in format 'ENV_KEY=secret-name/key'" {
		t.Fatalf("Expected error to happen, but was '%s'", err)
	}
}

func TestServiceSpecWithInvalidEnvConfigMap(t *testing.T) {
	serviceFlags := cmdflags.ServiceFlags{
		NamespaceFlags: cmdcore.NamespaceFlags{
			Name: "test-namespace",
		},
		Name: "test-service",
	}

	deployFlags := DeployFlags{
		Image:         "test-image",
		EnvConfigMaps: []string{"test-env-config-map-key1"},

		RemoveKnctlDeployEnvVar: true,
	}

	_, err := NewServiceSpec(serviceFlags, deployFlags).Service()
	if err == nil {
		t.Fatalf("Expected error to happen")
	}

	if err.Error() != "Expected environment variable from config map to be in format 'ENV_KEY=config-map-name/key'" {
		t.Fatalf("Expected error to happen, but was '%s'", err)
	}
}
