// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

const (
	IgnoreUnderlyingDeletionAnnotation = "nauticus.io/ignore-underlying-deletion"
)

func (s *Space) HasIgnoreUnderlyingDeletionAnnotation() bool {
	if _, ok := s.Annotations[IgnoreUnderlyingDeletionAnnotation]; ok {
		return true
	}

	return false
}
