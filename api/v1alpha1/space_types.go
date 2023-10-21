// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	ConditionType    string
	ConditionReason  string
	ConditionMessage string
)

// SpaceSpec defines the desired state of Space.
type SpaceSpec struct {
	// Specifies a list of ResourceQuota resources assigned to the Space. The assigned values are inherited by the namespace created by the Space. Optional.
	ResourceQuota corev1.ResourceQuotaSpec `json:"resourceQuota,omitempty"`
	// Specifies the owners of the Space. Mandatory.
	Owners []v1.Subject `json:"owners,omitempty"`
	// Specifies additional RoleBindings assigned to the Space. Nauticus will ensure that the namespace in the Space always contain the RoleBinding for the given ClusterRole. Optional.
	AdditionalRoleBindings AdditionalRoleBindingsSpec `json:"additionalRoleBindings,omitempty"`
	// Specifies the NetworkPolicies assigned to the Tenant. The assigned NetworkPolicies are inherited by the namespace created in the Space. Optional.
	NetworkPolicies NetworkPolicies `json:"networkPolicies,omitempty"`
	// Specifies the resource min/max usage restrictions to the Space. Optional.
	LimitRanges LimitRangesSpec `json:"limitRanges,omitempty"`
	// Specifies a list of service account to create within the Space. Optional
	ServiceAccounts ServiceAccountsSpec `json:"serviceAccounts,omitempty"`
	// Reference to a SpaceTemplate
	TemplateRef SpaceTemplateReference `json:"templateRef,omitempty"`
}

// SpaceTemplateReference.
type SpaceTemplateReference struct {
	// Name of the SpaceTemplate.
	Name string `json:"name,omitempty"`
	// Kind specifies the kind of the referenced resource, which should be "SpaceTemplate".
	Kind string `json:"kind,omitempty"`
	// Group is the API group of the SpaceTemplate,  "nauticus.io/v1alpha1".
	Group string `json:"group,omitempty"`
}

// SpaceStatus defines the observed state of Space.
type SpaceStatus struct {
	// NamespaceName the name of the created underlying namespace.
	NamespaceName string `json:"namespaceName,omitempty"`
	// Conditions List of status conditions to indicate the status of Space
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,categories={spaces}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",description="Ready"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Age"
// +kubebuilder:printcolumn:name="NamespaceName",type=string,JSONPath=`.status.namespaceName`

// Space is the Schema for the spaces API.
type Space struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceSpec   `json:"spec,omitempty"`
	Status SpaceStatus `json:"status,omitempty"`
}

func (s *Space) GetConditions() []metav1.Condition {
	return s.Status.Conditions
}

func (s *Space) SetConditions(conditions []metav1.Condition) {
	s.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// SpaceList contains a list of Space.
type SpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Space `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Space{}, &SpaceList{})
}
