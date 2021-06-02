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
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PubSubSubscriptionTemplateStatus defines the observed state of PubSubSubscriptionTemplate
type PubSubSubscriptionTemplateStatus struct {
	Ref v1.ObjectReference `json:"ref,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PubSubSubscriptionTemplate is the Schema for the pubsubsubscriptiontemplates API
type PubSubSubscriptionTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   v1beta1.PubSubSubscriptionSpec   `json:"spec,omitempty"`
	Status PubSubSubscriptionTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PubSubSubscriptionTemplateList contains a list of PubSubSubscriptionTemplate
type PubSubSubscriptionTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PubSubSubscriptionTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PubSubSubscriptionTemplate{}, &PubSubSubscriptionTemplateList{})
}