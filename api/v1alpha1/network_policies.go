package v1alpha1

import (
	networkingv1 "k8s.io/api/networking/v1"
)

type NetworkPolicies struct {
	EnableDefaultStrictMode bool                             `json:"enableDefaultStrictMode,omitempty"`
	Items                   []networkingv1.NetworkPolicySpec `json:"items,omitempty"`
}
