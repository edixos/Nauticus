// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import v1 "k8s.io/api/rbac/v1"

type AdditionalRoleBinding struct {
	RoleRef  v1.RoleRef   `json:"roleRef,omitempty"`
	Subjects []v1.Subject `json:"subjects,omitempty"`
}

type AdditionalRoleBindingsSpec []AdditionalRoleBinding
