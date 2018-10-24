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

package service

import (
	"fmt"
	"strings"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apirand "k8s.io/apimachinery/pkg/util/rand"
)

type ServiceSpec struct{}

func (ServiceSpec) Build(serviceFlags cmdflags.ServiceFlags, deployFlags DeployFlags) (v1alpha1.Service, error) {
	var buildSpec *buildv1alpha1.BuildSpec

	if deployFlags.BuildCreateArgsFlags.IsProvided() {
		// TODO assumes that same image is used for building and running
		deployFlags.BuildCreateArgsFlags.Image = deployFlags.Image

		spec, err := ctlbuild.BuildSpec{}.Build(deployFlags.BuildCreateArgsFlags.BuildSpecOpts)
		if err != nil {
			return v1alpha1.Service{}, err
		}

		buildSpec = &spec
	}

	serviceCont := corev1.Container{
		Image: deployFlags.Image,
	}

	for _, kv := range deployFlags.Env {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return v1alpha1.Service{}, fmt.Errorf("Expected environment variable to be in format 'KEY=VALUE'")
		}
		serviceCont.Env = append(serviceCont.Env, corev1.EnvVar{Name: pieces[0], Value: pieces[1]})
	}

	// TODO it's convenient to force redeploy anytime deploy is issued
	if !deployFlags.RemoveKnctlDeployEnvVar {
		serviceCont.Env = append(serviceCont.Env, corev1.EnvVar{
			Name:  "KNCTL_DEPLOY",
			Value: apirand.String(10),
		})
	}

	service := v1alpha1.Service{
		ObjectMeta: deployFlags.GenerateNameFlags.Apply(metav1.ObjectMeta{
			Name:      serviceFlags.Name,
			Namespace: serviceFlags.NamespaceFlags.Name,
		}),
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					Build: buildSpec,
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							// TODO service account may be different for runtime vs build
							ServiceAccountName: deployFlags.BuildCreateArgsFlags.ServiceAccountName,
							Container:          serviceCont,
						},
					},
				},
			},
		},
	}

	return service, nil
}
