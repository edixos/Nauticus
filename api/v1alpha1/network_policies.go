// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	networkingv1 "k8s.io/api/networking/v1"
)

type NetworkPolicies struct {
	EnableDefaultStrictMode bool                             `json:"enableDefaultStrictMode,omitempty"`
	Items                   []networkingv1.NetworkPolicySpec `json:"items,omitempty"`
}
