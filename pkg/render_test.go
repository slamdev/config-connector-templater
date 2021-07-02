/*
Copyright 2021.

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

package pkg

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	api "github.com/slamdev/config-connector-templater/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestRender(t *testing.T) {
	stringPtr := func(s string) *string {
		return &s
	}

	template := api.PubSubTopicTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
		},
		Spec: pubsub.PubSubTopicSpec{
			KmsKeyRef: &v1alpha1.ResourceRef{
				Name:      "{{ .metadata.namespace }}-ref-name",
				Namespace: "ref-ns",
			},
			MessageStoragePolicy: &pubsub.TopicMessageStoragePolicy{
				AllowedPersistenceRegions: []string{"{{ .metadata.namespace }}-r1"},
			},
			ResourceID: stringPtr("{{ .metadata.namespace }}.test-ref"),
		},
	}

	res, err := Render(template.Spec, template)
	if err != nil {
		t.Error(err)
	}
	resSpec := res.(pubsub.PubSubTopicSpec)

	assert.Equal(t, "test-ns.test-ref", *resSpec.ResourceID)
	assert.Equal(t, "test-ns-r1", resSpec.MessageStoragePolicy.AllowedPersistenceRegions[0])
	assert.Equal(t, "test-ns-ref-name", resSpec.KmsKeyRef.Name)
	assert.Equal(t, template.Spec.KmsKeyRef.Namespace, resSpec.KmsKeyRef.Namespace)
}
