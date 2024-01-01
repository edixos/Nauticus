// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

type ServiceAccountsSpec struct {
	// Specifies the list of Service Account to be created. Optional
	Items []ServiceAccountSpec `json:"items,omitempty"`
}

type Annotations map[string]string

type ServiceAccountSpec struct {
	// Specifies the service account name to be created. Required
	Name string `json:"name,omitempty"`
	// Specifies the annotations to be placed in the ServiceAccount. Optional
	Annotations Annotations `json:"annotations,omitempty"`
}
