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

package v1alpha1

import (
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PubSubTopicTemplateStatus defines the observed state of PubSubTopicTemplate
type PubSubTopicTemplateStatus struct {
	Ref v1.ObjectReference `json:"ref,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PubSubTopicTemplate is the Schema for the pubsubtopictemplates API
type PubSubTopicTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   pubsub.PubSubTopicSpec    `json:"spec,omitempty"`
	Status PubSubTopicTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PubSubTopicTemplateList contains a list of PubSubTopicTemplate
type PubSubTopicTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PubSubTopicTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PubSubTopicTemplate{}, &PubSubTopicTemplateList{})
}
