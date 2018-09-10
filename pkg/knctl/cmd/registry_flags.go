/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Open 2.0 (the "License");
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
	"github.com/cppforlife/knctl/pkg/knctl/registry"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type RegistryFlags struct {
	ClusterRegistry          bool
	ClusterRegistryNamespace string
}

func (s *RegistryFlags) Set(cmd *cobra.Command, flagsFactory FlagsFactory) {
	cmd.Flags().BoolVarP(&s.ClusterRegistry, "cluster-registry", "c", false, "Use cluster registry")
	cmd.Flags().StringVar(&s.ClusterRegistryNamespace, "cluster-registry-namespace", "", "Namespace where cluster registry was installed")
}

type RegistryConfiguration struct {
	rf         RegistryFlags
	coreClient kubernetes.Interface
}

func NewRegistryConfiguration(rf RegistryFlags, coreClient kubernetes.Interface) RegistryConfiguration {
	return RegistryConfiguration{rf, coreClient}
}

func (r RegistryConfiguration) ApplyToService(sf ServiceFlags, df *DeployFlags) error {
	info, err := r.imageInfo(sf.NamespaceFlags)
	if err != nil {
		return err
	}

	if info != nil {
		// TODO generate name???
		df.Image = info.Address + "/" + sf.NamespaceFlags.Name + "/" + sf.Name
		df.BuildCreateArgsFlags.ServiceAccountName = info.ServiceAccountName
		r.addCACerts(&df.BuildCreateArgsFlags)
	}

	return nil
}

func (r RegistryConfiguration) ApplyToBuild(bf BuildFlags, bcf *BuildCreateFlags) error {
	info, err := r.imageInfo(bf.NamespaceFlags)
	if err != nil {
		return err
	}

	if info != nil {
		// TODO generate name???
		bcf.Image = info.Address + "/" + bf.NamespaceFlags.Name + "/" + bf.Name
		bcf.ServiceAccountName = info.ServiceAccountName
		r.addCACerts(&bcf.BuildCreateArgsFlags)
	}

	return nil
}

func (c RegistryConfiguration) addCACerts(bcf *BuildCreateArgsFlags) {
	const certsDir = "/etc/ssl/certs/"

	bcf.VolumeMounts = []corev1.VolumeMount{
		{Name: "host-certs", ReadOnly: true, MountPath: certsDir},
	}

	bcf.Volumes = []corev1.Volume{
		{
			Name: "host-certs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{Path: certsDir},
			},
		},
	}

	bcf.Env = []corev1.EnvVar{
		{Name: "SSL_CERT_DIR", Value: certsDir},
	}
}

func (r RegistryConfiguration) imageInfo(nf NamespaceFlags) (*registry.ImageInfo, error) {
	var registryInfoSrc registry.RegistryInfoSource

	switch {
	case r.rf.ClusterRegistry:
		registryInfoSrc = registry.NewClusterRegistry(r.rf.ClusterRegistryNamespace, r.coreClient)
	}

	if registryInfoSrc != nil {
		info, err := registry.NewRegistry(nf.Name, r.coreClient, registryInfoSrc).Info()
		if err != nil {
			return nil, err
		}
		return &info, nil
	}

	return nil, nil
}
