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
	"sort"

	ctlbuild "github.com/cppforlife/knctl/pkg/knctl/build"
	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	buildclientset "github.com/knative/build/pkg/client/clientset/versioned"
	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	service         v1alpha1.Service
	servingClient   servingclientset.Interface
	buildClient     buildclientset.Interface
	coreClient      kubernetes.Interface
	buildObjFactory ctlbuild.Factory
}

func NewService(
	service v1alpha1.Service,
	servingClient servingclientset.Interface,
	buildClient buildclientset.Interface,
	coreClient kubernetes.Interface,
	buildObjFactory ctlbuild.Factory,
) Service {
	return Service{service, servingClient, buildClient, coreClient, buildObjFactory}
}

func (l Service) CreatedBuildSinceRevision(lastRevision *v1alpha1.Revision) (ctlbuild.Build, error) {
	cancelResWatchCh := make(chan struct{})
	revisionsToWatchCh := make(chan v1alpha1.Revision)
	buildsToWatchCh := make(chan buildv1alpha1.Build)

	// Watch revisions in this service
	go func() {
		revisionWatcher := NewRevisionWatcher(
			l.servingClient.ServingV1alpha1().Revisions(l.service.Namespace),
			metav1.ListOptions{
				LabelSelector: labels.Set(map[string]string{
					serving.ConfigurationLabelKey: l.service.Name,
				}).String(),
			},
		)

		err := revisionWatcher.Watch(revisionsToWatchCh, cancelResWatchCh)
		if err != nil {
			// TODO error?
			fmt.Printf("Revision watching error: %s\n", err)
		}

		close(revisionsToWatchCh)
	}()

	var createdRevision *v1alpha1.Revision

	for revision := range revisionsToWatchCh {
		if lastRevision == nil || revision.CreationTimestamp.Time.After(lastRevision.CreationTimestamp.Time) {
			createdRevision = &revision
			break
		}
	}

	if createdRevision == nil {
		return ctlbuild.Build{}, fmt.Errorf("Expected to find created revision")
	}

	// TODO build is not not necessarily in same namespace as service
	buildsClient := l.buildClient.BuildV1alpha1().Builds(l.service.Namespace)

	// Watch builds for new revision
	go func() {
		buildWatcher := ctlbuild.NewBuildWatcher(
			buildsClient,
			metav1.ListOptions{
				FieldSelector: fields.OneTermEqualSelector("metadata.name", createdRevision.Name).String(),
			},
		)

		err := buildWatcher.Watch(buildsToWatchCh, cancelResWatchCh)
		if err != nil {
			// TODO error?
			fmt.Printf("Build watching error: %s\n", err)
		}

		close(buildsToWatchCh)
	}()

	for build := range buildsToWatchCh {
		if build.Name == createdRevision.Name {
			close(cancelResWatchCh)
			return l.buildObjFactory.New(&build), nil
		}
	}

	return ctlbuild.Build{}, fmt.Errorf("Expected to find new build")
}

func (l Service) CreatedRevisionSinceRevision(lastRevision *v1alpha1.Revision) (*v1alpha1.Revision, error) {
	cancelResWatchCh := make(chan struct{})
	revisionsToWatchCh := make(chan v1alpha1.Revision)

	// Watch revisions in this service
	go func() {
		revisionWatcher := NewRevisionWatcher(
			l.servingClient.ServingV1alpha1().Revisions(l.service.Namespace),
			metav1.ListOptions{
				LabelSelector: labels.Set(map[string]string{
					serving.ConfigurationLabelKey: l.service.Name,
				}).String(),
			},
		)

		err := revisionWatcher.Watch(revisionsToWatchCh, cancelResWatchCh)
		if err != nil {
			// TODO error?
			fmt.Printf("Revision watching error: %s\n", err)
		}

		close(revisionsToWatchCh)
	}()

	var createdRevision *v1alpha1.Revision

	for revision := range revisionsToWatchCh {
		if lastRevision == nil || revision.CreationTimestamp.Time.After(lastRevision.CreationTimestamp.Time) {
			createdRevision = &revision
			break
		}
	}

	if createdRevision == nil {
		return nil, fmt.Errorf("Expected to find created revision")
	}

	return createdRevision, nil
}

func (l Service) LastRevision() (*v1alpha1.Revision, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			serving.ConfigurationLabelKey: l.service.Name,
		}).String(),
	}

	// TODO LatestCreatedRevisionName may not be updated that quickly
	revisions, err := l.servingClient.ServingV1alpha1().Revisions(l.service.Namespace).List(listOpts)
	if err != nil {
		return nil, err
	}

	if len(revisions.Items) == 0 {
		return nil, nil
	}

	sort.Slice(revisions.Items, func(i, j int) bool {
		return revisions.Items[i].CreationTimestamp.Time.After(revisions.Items[j].CreationTimestamp.Time)
	})

	revision := revisions.Items[0]
	return &revision, nil
}
