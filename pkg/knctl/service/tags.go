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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/knative/serving/pkg/apis/serving"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

const (
	tagRevisionLabelKeyPrefix = "tag.cli.knative.dev/"
	tagRevisionLabelValue     = "true"

	TagsLatest   = "latest"
	TagsPrevious = "previous"
)

var (
	rfc6901Encoder = strings.NewReplacer("~", "~0", "/", "~1", ":", "~7")
)

type Tags struct {
	servingClient servingclientset.Interface
}

func NewTags(servingClient servingclientset.Interface) Tags {
	return Tags{servingClient}
}

func (t Tags) List(revision v1alpha1.Revision) []string {
	var tags []string
	for k, _ := range revision.Labels {
		if strings.HasPrefix(k, tagRevisionLabelKeyPrefix) {
			tags = append(tags, strings.TrimPrefix(k, tagRevisionLabelKeyPrefix))
		}
	}
	return tags
}

type TagsFindContext struct {
	Namespace string
	Service   string
}

// TODO instead of context, pass in v1alpha1.Service?
func (t Tags) Find(context TagsFindContext, tag string) (*v1alpha1.Revision, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			serving.ConfigurationLabelKey: context.Service,
			t.label(tag):                  tagRevisionLabelValue,
		}).String(),
	}

	revisions, err := t.servingClient.ServingV1alpha1().Revisions(context.Namespace).List(listOpts)
	if err != nil {
		return nil, fmt.Errorf("Listing revisions: %s", err)
	}

	switch len(revisions.Items) {
	case 0:
		return nil, fmt.Errorf("Expected to find revision with tag '%s' for service '%s'", tag, context.Service)

	case 1:
		revision := revisions.Items[0]
		return &revision, nil

	default:
		return nil, fmt.Errorf("Expected to find extactly one revision with tag '%s' for service '%s', but found multiple", tag, context.Service)
	}
}

func (t Tags) Repoint(revision *v1alpha1.Revision, tag string) error {
	listOpts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			serving.ConfigurationLabelKey: revision.Labels[serving.ConfigurationLabelKey],
			t.label(tag):                  tagRevisionLabelValue,
		}).String(),
	}

	revisions, err := t.servingClient.ServingV1alpha1().Revisions(revision.Namespace).List(listOpts)
	if err != nil {
		return fmt.Errorf("Listing revisions: %s", err)
	}

	for _, revision := range revisions.Items {
		err := t.Untag(revision, tag)
		if err != nil {
			return err
		}
	}

	return t.Tag(revision, tag)
}

func (t Tags) Tag(revision *v1alpha1.Revision, tag string) error {
	// TODO support multiple tags
	// TODO better namespacing?
	mergePatch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				t.label(tag): tagRevisionLabelValue,
			},
		},
	}

	patchJSON, err := json.Marshal(mergePatch)
	if err != nil {
		return err
	}

	var lastErr error

	// TODO better way to avoid race: 'unable to find api field in struct...'
	for i := 0; i < 5; i++ {
		_, lastErr = t.servingClient.ServingV1alpha1().Revisions(revision.Namespace).Patch(revision.Name, types.MergePatchType, patchJSON)
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("Tagging revision: %s", lastErr)
}

func (t Tags) Untag(revision v1alpha1.Revision, tag string) error {
	encodedTag := rfc6901Encoder.Replace(t.label(tag))
	patchJSON := []byte(fmt.Sprintf(`[{"op":"remove", "path":"/metadata/labels/%s"}]`, encodedTag))

	var lastErr error

	// TODO better way to avoid race: 'unable to find api field in struct...'
	for i := 0; i < 5; i++ {
		_, lastErr = t.servingClient.ServingV1alpha1().Revisions(revision.Namespace).Patch(revision.Name, types.JSONPatchType, patchJSON)
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("Untagging revision: %s", lastErr)
}

func (t Tags) label(tag string) string {
	return tagRevisionLabelKeyPrefix + tag
}
