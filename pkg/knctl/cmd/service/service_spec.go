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
	"strconv"
	"strings"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	cmdflags "github.com/cppforlife/knctl/pkg/knctl/cmd/flags"
	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apirand "k8s.io/apimachinery/pkg/util/rand"
)

type ServiceSpec struct {
	serviceFlags cmdflags.ServiceFlags
	deployFlags  DeployFlags
}

func NewServiceSpec(serviceFlags cmdflags.ServiceFlags, deployFlags DeployFlags) ServiceSpec {
	return ServiceSpec{serviceFlags, deployFlags}
}

func (s ServiceSpec) Namespace() string { return s.serviceFlags.NamespaceFlags.Name }
func (s ServiceSpec) Name() string      { return s.serviceFlags.Name }

func (s ServiceSpec) HasBuild() bool {
	return s.deployFlags.BuildCreateArgsFlags.IsProvided()
}

func (s ServiceSpec) NeedsConfigurationUpdate() bool {
	return !s.deployFlags.ManagedRoute
}

func (s ServiceSpec) Service() (v1alpha1.Service, error) {
	service := v1alpha1.Service{
		ObjectMeta: s.deployFlags.GenerateNameFlags.Apply(metav1.ObjectMeta{
			Name:      s.serviceFlags.Name,
			Namespace: s.serviceFlags.NamespaceFlags.Name,
		}),
	}

	if s.NeedsConfigurationUpdate() {
		service.Spec.Manual = &v1alpha1.ManualType{}
	} else {
		conf, err := s.Configuration()
		if err != nil {
			return v1alpha1.Service{}, err
		}

		service.Spec.RunLatest = &v1alpha1.RunLatestType{
			Configuration: conf.Spec,
		}
	}

	return service, nil
}

func (s ServiceSpec) Configuration() (v1alpha1.Configuration, error) {
	var buildSpec *buildv1alpha1.BuildSpec

	if s.deployFlags.BuildCreateArgsFlags.IsProvided() {
		// TODO assumes that same image is used for building and running
		s.deployFlags.BuildCreateArgsFlags.Image = s.deployFlags.Image

		spec, err := ctlbuild.BuildSpec{}.Build(s.deployFlags.BuildCreateArgsFlags.BuildSpecOpts)
		if err != nil {
			return v1alpha1.Configuration{}, err
		}

		buildSpec = &spec
	}

	serviceCont := corev1.Container{
		Image: s.deployFlags.Image,
	}

	for _, kv := range s.deployFlags.EnvVars {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return v1alpha1.Configuration{}, fmt.Errorf("Expected environment variable to be in format 'ENV_KEY=value'")
		}
		serviceCont.Env = append(serviceCont.Env, corev1.EnvVar{Name: pieces[0], Value: pieces[1]})
	}

	envVars, err := s.buildEnvFromSecrets(s.deployFlags)
	if err != nil {
		return v1alpha1.Configuration{}, err
	}

	serviceCont.Env = append(serviceCont.Env, envVars...)

	envVars, err = s.buildEnvFromConfigMaps(s.deployFlags)
	if err != nil {
		return v1alpha1.Configuration{}, err
	}

	serviceCont.Env = append(serviceCont.Env, envVars...)

	for _, configmapName := range s.deployFlags.EnvFromConfigMaps {
		envFromSource := corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource {
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configmapName,
				},
			},
		}

		serviceCont.EnvFrom = append(serviceCont.EnvFrom, envFromSource)
	}

	secretVolumes, secretVolumeMounts, err := s.buildMountToSecrets(s.deployFlags)
	if err != nil {
		return v1alpha1.Configuration{}, err
	}

	configMapVolumes, configMapVolumeMounts, err:= s.buildMountToConfigMaps(s.deployFlags)
	if err != nil {
		return v1alpha1.Configuration{}, err
	}

	serviceCont.VolumeMounts = append(serviceCont.VolumeMounts, secretVolumeMounts...)
	serviceCont.VolumeMounts = append(serviceCont.VolumeMounts, configMapVolumeMounts...)

	// TODO it's convenient to force redeploy anytime deploy is issued
	if !s.deployFlags.RemoveKnctlDeployEnvVar {
		serviceCont.Env = append(serviceCont.Env, corev1.EnvVar{
			Name:  "KNCTL_DEPLOY",
			Value: apirand.String(10),
		})
	}

	revisionAnns := map[string]string{}

	if s.deployFlags.MinScale != nil {
		revisionAnns["autoscaling.knative.dev/minScale"] = strconv.Itoa(*s.deployFlags.MinScale)
	}
	if s.deployFlags.MaxScale != nil {
		revisionAnns["autoscaling.knative.dev/maxScale"] = strconv.Itoa(*s.deployFlags.MaxScale)
	}

	conf := v1alpha1.Configuration{
		// ObjectMeta is populated when object is being created
		Spec: v1alpha1.ConfigurationSpec{
			Build: &v1alpha1.RawExtension{BuildSpec: buildSpec},
			RevisionTemplate: v1alpha1.RevisionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: revisionAnns,
				},
				Spec: v1alpha1.RevisionSpec{
					// TODO service account may be different for runtime vs build
					ServiceAccountName: s.deployFlags.BuildCreateArgsFlags.ServiceAccountName,
					Container:          serviceCont,
				},
			},
		},
	}

	if len(secretVolumes) > 0 {
		conf.Spec.RevisionTemplate.Spec.Volumes = append(conf.Spec.RevisionTemplate.Spec.Volumes, secretVolumes...)
	}

	if len(configMapVolumes) > 0 {
		conf.Spec.RevisionTemplate.Spec.Volumes = append(conf.Spec.RevisionTemplate.Spec.Volumes, configMapVolumes...)
	}

	if s.deployFlags.ContainerConcurrency != nil {
		conf.Spec.RevisionTemplate.Spec.ContainerConcurrency = v1alpha1.RevisionContainerConcurrencyType(*s.deployFlags.ContainerConcurrency)
	}

	return conf, nil
}

func (s ServiceSpec) buildEnvFromSecrets(deployFlags DeployFlags) ([]corev1.EnvVar, error) {
	var result []corev1.EnvVar

	for _, kv := range s.deployFlags.EnvSecrets {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected environment variable from secret to be in format 'ENV_KEY=secret-name/key'")
		}

		secretPieces := strings.SplitN(pieces[1], "/", 2)
		if len(secretPieces) != 2 {
			return nil, fmt.Errorf("Expected environment variable secret ref to be in format 'secret-name/key'")
		}

		result = append(result, corev1.EnvVar{
			Name: pieces[0],
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretPieces[0],
					},
					Key: secretPieces[1],
				},
			},
		})
	}

	return result, nil
}

func (s ServiceSpec) buildEnvFromConfigMaps(deployFlags DeployFlags) ([]corev1.EnvVar, error) {
	var result []corev1.EnvVar

	for _, kv := range s.deployFlags.EnvConfigMaps {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Expected environment variable from config map to be in format 'ENV_KEY=config-map-name/key'")
		}

		mapPieces := strings.SplitN(pieces[1], "/", 2)
		if len(mapPieces) != 2 {
			return nil, fmt.Errorf("Expected environment variable config map ref to be in format 'config-map-name/key'")
		}

		result = append(result, corev1.EnvVar{
			Name: pieces[0],
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mapPieces[0],
					},
					Key: mapPieces[1],
				},
			},
		})
	}

	return result, nil
}


func (s ServiceSpec) buildMountToSecrets(deployFlags DeployFlags) ([]corev1.Volume, []corev1.VolumeMount, error) {
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	for _, kv := range s.deployFlags.VolumeMountSecrets {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, nil, fmt.Errorf("Expected parameter for mounting a secret should be in format 'secret-name=/mount/path'")
		}

		volumeName := fmt.Sprintf("secret-volume-%s", pieces[0])
		
		volumes = append(volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: pieces[0],
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name: volumeName,
			ReadOnly: true,
			MountPath: pieces[1],
		})
	}

	return volumes, volumeMounts, nil
}

func (s ServiceSpec) buildMountToConfigMaps(deployFlags DeployFlags) ([]corev1.Volume, []corev1.VolumeMount, error) {
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	for _, kv := range s.deployFlags.VolumeMountConfigMaps {
		pieces := strings.SplitN(kv, "=", 2)
		if len(pieces) != 2 {
			return nil, nil, fmt.Errorf("Expected parameter for mounting a config map should be in format 'config-map-name=/mount/path'")
		}

		volumeName := fmt.Sprintf("configmap-volume-%s", pieces[0])
		
		volumes = append(volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: pieces[0],
					},
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name: volumeName,
			ReadOnly: true,
			MountPath: pieces[1],
		})
	}

	return volumes, volumeMounts, nil
}