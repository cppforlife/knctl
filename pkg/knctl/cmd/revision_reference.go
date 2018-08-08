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
	"fmt"
	"strings"

	ctlservice "github.com/cppforlife/knctl/pkg/knctl/service"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RevisionReference struct {
	flags         RevisionFlags
	tags          ctlservice.Tags
	servingClient servingclientset.Interface
}

func NewRevisionReference(flags RevisionFlags, tags ctlservice.Tags, servingClient servingclientset.Interface) RevisionReference {
	return RevisionReference{flags, tags, servingClient}
}

func (r RevisionReference) Revision() (*v1alpha1.Revision, error) {
	if len(r.flags.Name) == 0 {
		return nil, fmt.Errorf("Expected revision reference to not be empty")
	}

	pieces := strings.Split(r.flags.Name, ":")

	switch len(pieces) {
	case 1:
		revision, err := r.servingClient.ServingV1alpha1().Revisions(r.flags.NamespaceFlags.Name).Get(pieces[0], metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Getting revision: %s", err)
		}
		return revision, err

	case 2:
		context := ctlservice.TagsFindContext{Namespace: r.flags.NamespaceFlags.Name, Service: pieces[0]}
		return r.tags.Find(context, pieces[1])

	default:
		return nil, fmt.Errorf("Expected revision reference to be in format 'revision' or 'service:revision'")
	}
}
