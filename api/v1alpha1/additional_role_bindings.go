package v1alpha1

import v1 "k8s.io/api/rbac/v1"

type AdditionalRoleBinding struct {
    ClusterRoleName string       `json:"clusterRoleName,omitempty"`
    Subjects        []v1.Subject `json:"subjects,omitempty"`
}

type AdditionalRoleBindings []AdditionalRoleBinding
