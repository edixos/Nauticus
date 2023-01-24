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
    v1 "k8s.io/api/rbac/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SpaceSpec defines the desired state of Space
type SpaceSpec struct {
    // Specifies a list of ResourceQuota resources assigned to the Space. The assigned values are inherited by the namespace created by the Space. Optional
    ResourceQuota corev1.ResourceQuotaSpec `json:"resourceQuota,omitempty"`
    // Specifies the owners of the Space. Mandatory
    Owners []v1.Subject `json:"owners,omitempty"`
    // Specifies additional RoleBindings assigned to the Space. Nauticus will ensure that the namespace in the Space always contain the RoleBinding for the given ClusterRole. Optional
    AdditionalRoleBindings AdditionalRoleBindings `json:"additionalRoleBindings,omitempty"`
    // Specifies the NetworkPolicies assigned to the Tenant. The assigned NetworkPolicies are inherited by the namespace created in the Space. Optional.
    NetworkPolicies NetworkPolicies `json:"networkPolicies,omitempty"`
}

// SpaceStatus defines the observed state of Space
type SpaceStatus struct {
    // NamespaceName the name of the created underlying namespace.
    NamespaceName string `json:"namespaceName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,categories={spaces}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="NamespaceName",type=string,JSONPath=`.status.namespaceName`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Age"

// Space is the Schema for the spaces API
type Space struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   SpaceSpec   `json:"spec,omitempty"`
    Status SpaceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpaceList contains a list of Space
type SpaceList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Space `json:"items"`
}

func init() {
    SchemeBuilder.Register(&Space{}, &SpaceList{})
}
