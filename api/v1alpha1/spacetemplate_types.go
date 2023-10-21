// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpaceTemplateSpec defines the desired state of SpaceTemplate.
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

// SpaceTemplateStatus defines the observed state of SpaceTemplate.
type SpaceTemplateStatus struct {
	// Status is the status of the cluster.
	Status string `json:"status"`
	// Conditions List of status conditions to indicate the status of Space
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster, categories={spacetemplate}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",description="Ready"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",description="Ready"

// SpaceTemplate is the Schema for the spacetemplates API.
type SpaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceTemplateSpec   `json:"spec,omitempty"`
	Status SpaceTemplateStatus `json:"status,omitempty"`
}

func (in *SpaceTemplate) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

func (in *SpaceTemplate) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// SpaceTemplateList contains a list of SpaceTemplate.
type SpaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpaceTemplate{}, &SpaceTemplateList{})
}
