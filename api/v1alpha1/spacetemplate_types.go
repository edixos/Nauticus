/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpaceTemplateSpec defines the desired state of SpaceTemplate
type SpaceTemplateSpec struct {
	// Specifies a list of ResourceQuota resources assigned to the Space. The assigned values are inherited by the namespace created by the Space. Optional.
	ResourceQuota corev1.ResourceQuotaSpec `json:"resourceQuota,omitempty"`
	// Specifies additional RoleBindings assigned to the Space. Nauticus will ensure that the namespace in the Space always contain the RoleBinding for the given ClusterRole. Optional.
	AdditionalRoleBindings AdditionalRoleBindingsSpec `json:"additionalRoleBindings,omitempty"`
	// Specifies the NetworkPolicies assigned to the Tenant. The assigned NetworkPolicies are inherited by the namespace created in the Space. Optional.
	NetworkPolicies NetworkPolicies `json:"networkPolicies,omitempty"`
	// Specifies the resource min/max usage restrictions to the Space. Optional.
	LimitRanges LimitRangesSpec `json:"limitRanges,omitempty"`
}

// SpaceTemplateStatus defines the observed state of SpaceTemplate
type SpaceTemplateStatus struct {
	// Conditions List of status conditions to indicate the status of Space
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SpaceTemplate is the Schema for the spacetemplates API
type SpaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceTemplateSpec   `json:"spec,omitempty"`
	Status SpaceTemplateStatus `json:"status,omitempty"`
}

//func (st *SpaceTemplate) ApplyToDelete(options *client.DeleteOptions) {
//TODO implement me
//	panic("implement me")
//}

//+kubebuilder:object:root=true

// SpaceTemplateList contains a list of SpaceTemplate
type SpaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpaceTemplate{}, &SpaceTemplateList{})
}
