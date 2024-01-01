// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0
package space

import (
	"errors"
	"reflect"
	"testing"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestMergeResourceQuotas(t *testing.T) {
	testCases := []struct {
		name          string
		space         *nauticusiov1alpha1.Space
		spaceTemplate *nauticusiov1alpha1.SpaceTemplate
		expected      *corev1.ResourceQuotaSpec
		expectedErr   error
	}{
		{
			name: "Both space and spaceTemplate ResourceQuotas provided",
			space: &nauticusiov1alpha1.Space{
				Spec: nauticusiov1alpha1.SpaceSpec{
					ResourceQuota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceLimitsCPU:      resource.MustParse("8"),
							corev1.ResourceLimitsMemory:   resource.MustParse("16Gi"),
							corev1.ResourceRequestsCPU:    resource.MustParse("4"),
							corev1.ResourceRequestsMemory: resource.MustParse("8Gi"),
						},
					},
				},
			},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
				Spec: nauticusiov1alpha1.SpaceTemplateSpec{
					ResourceQuota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceLimitsCPU:      resource.MustParse("2"),
							corev1.ResourceLimitsMemory:   resource.MustParse("2Gi"),
							corev1.ResourceRequestsCPU:    resource.MustParse("1"),
							corev1.ResourceRequestsMemory: resource.MustParse("1Gi"),
						},
					},
				},
			},
			expected: &corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceLimitsCPU:      resource.MustParse("8"),
					corev1.ResourceLimitsMemory:   resource.MustParse("16Gi"),
					corev1.ResourceRequestsCPU:    resource.MustParse("4"),
					corev1.ResourceRequestsMemory: resource.MustParse("8Gi"),
				},
			},
			expectedErr: nil,
		},
		{
			name: "Both space and spaceTemplate ResourceQuotas provided (CPU)",
			space: &nauticusiov1alpha1.Space{
				Spec: nauticusiov1alpha1.SpaceSpec{
					ResourceQuota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceLimitsCPU:      resource.MustParse("8"),
							corev1.ResourceRequestsCPU:    resource.MustParse("4"),
							corev1.ResourceLimitsMemory:   resource.MustParse("1Gi"),
							corev1.ResourceRequestsMemory: resource.MustParse("500Mi"),
						},
					},
				},
			},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
				Spec: nauticusiov1alpha1.SpaceTemplateSpec{
					ResourceQuota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceLimitsCPU:      resource.MustParse("2"),
							corev1.ResourceRequestsCPU:    resource.MustParse("1"),
							corev1.ResourceLimitsMemory:   resource.MustParse("1Gi"),
							corev1.ResourceRequestsMemory: resource.MustParse("500Mi"),
						},
					},
				},
			},
			expected: &corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceLimitsCPU:      resource.MustParse("8"),
					corev1.ResourceRequestsCPU:    resource.MustParse("4"),
					corev1.ResourceLimitsMemory:   resource.MustParse("1Gi"),
					corev1.ResourceRequestsMemory: resource.MustParse("500Mi"),
				},
			},
			expectedErr: nil,
		},
		{
			name:  "Only spaceTemplate ResourceQuotas (limits) provided",
			space: &nauticusiov1alpha1.Space{},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
				Spec: nauticusiov1alpha1.SpaceTemplateSpec{
					ResourceQuota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceLimitsCPU:    resource.MustParse("3"),
							corev1.ResourceLimitsMemory: resource.MustParse("3Gi"),
						},
					},
				},
			},
			expected: &corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceLimitsCPU:    resource.MustParse("3"),
					corev1.ResourceLimitsMemory: resource.MustParse("3Gi"),
				},
			},
			expectedErr: nil,
		},
		{
			name:          "Both space and spaceTemplate ResourceQuotas are empty",
			space:         &nauticusiov1alpha1.Space{},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{},
			expected:      nil,
			expectedErr:   errors.New("merge not required both space and spacetpl resource quotas are empty"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MergeResourceQuotas(tc.space, tc.spaceTemplate)

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("Expected: %v, Got: %v", tc.expected, result)
			}
			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Errorf("Expected error: %v, Got error: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMergeRoleBindings(t *testing.T) {
	testCases := []struct {
		name          string
		space         *nauticusiov1alpha1.Space
		spaceTemplate *nauticusiov1alpha1.SpaceTemplate
		expected      []nauticusiov1alpha1.AdditionalRoleBinding
		expectedErr   error
	}{
		{
			name: "Space has role bindings, SpaceTemplate has role bindings",
			space: &nauticusiov1alpha1.Space{
				Spec: nauticusiov1alpha1.SpaceSpec{
					AdditionalRoleBindings: []nauticusiov1alpha1.AdditionalRoleBinding{
						{
							RoleRef: v1.RoleRef{
								APIGroup: "rbac.authorization.k8s.io",
								Kind:     "ClusterRole",
								Name:     "editor",
							},
							Subjects: []v1.Subject{
								{
									Name: "bob",
									Kind: "User",
								},
								{
									Name: "dev2",
									Kind: "Group",
								},
							},
						},
					},
				},
			},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
				Spec: nauticusiov1alpha1.SpaceTemplateSpec{
					AdditionalRoleBindings: []nauticusiov1alpha1.AdditionalRoleBinding{
						{
							RoleRef: v1.RoleRef{
								APIGroup: "rbac.authorization.k8s.io",
								Kind:     "ClusterRole",
								Name:     "viewer",
							},
							Subjects: []v1.Subject{
								{
									Name: "alice",
									Kind: "User",
								},
								{
									Name: "dev",
									Kind: "Group",
								},
							},
						},
					},
				},
			},
			expected: []nauticusiov1alpha1.AdditionalRoleBinding{
				{
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "editor",
					},
					Subjects: []v1.Subject{
						{
							Name: "bob",
							Kind: "User",
						},
						{
							Name: "dev2",
							Kind: "Group",
						},
					},
				},
				{
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "viewer",
					},
					Subjects: []v1.Subject{
						{
							Name: "alice",
							Kind: "User",
						},
						{
							Name: "dev",
							Kind: "Group",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "Space has no role bindings, SpaceTemplate has role bindings",
			space: &nauticusiov1alpha1.Space{},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
				Spec: nauticusiov1alpha1.SpaceTemplateSpec{
					AdditionalRoleBindings: []nauticusiov1alpha1.AdditionalRoleBinding{
						{
							RoleRef: v1.RoleRef{
								APIGroup: "rbac.authorization.k8s.io",
								Kind:     "ClusterRole",
								Name:     "viewer",
							},
							Subjects: []v1.Subject{
								{
									Name: "alice",
									Kind: "User",
								},
								{
									Name: "dev",
									Kind: "Group",
								},
							},
						},
					},
				},
			},
			expected: []nauticusiov1alpha1.AdditionalRoleBinding{
				{
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "viewer",
					},
					Subjects: []v1.Subject{
						{
							Name: "alice",
							Kind: "User",
						},
						{
							Name: "dev",
							Kind: "Group",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Space has role bindings, SpaceTemplate has no role bindings",
			space: &nauticusiov1alpha1.Space{
				Spec: nauticusiov1alpha1.SpaceSpec{
					AdditionalRoleBindings: []nauticusiov1alpha1.AdditionalRoleBinding{
						{
							RoleRef: v1.RoleRef{
								APIGroup: "rbac.authorization.k8s.io",
								Kind:     "ClusterRole",
								Name:     "viewer",
							},
							Subjects: []v1.Subject{
								{
									Name: "alice",
									Kind: "User",
								},
								{
									Name: "dev",
									Kind: "Group",
								},
							},
						},
					},
				},
			},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{},
			expected: []nauticusiov1alpha1.AdditionalRoleBinding{
				{
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "viewer",
					},
					Subjects: []v1.Subject{
						{
							Name: "alice",
							Kind: "User",
						},
						{
							Name: "dev",
							Kind: "Group",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:          "Both Space and SpaceTemplate have no role bindings",
			space:         &nauticusiov1alpha1.Space{},
			spaceTemplate: &nauticusiov1alpha1.SpaceTemplate{},
			expected:      nil,
			expectedErr:   errors.New("no additional roles bindings merged from the template"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MergeRoleBindings(tc.space, tc.spaceTemplate)

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("Expected: %v, Got: %v", tc.expected, result)
			}
			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Errorf("Expected error: %v, Got error: %v", tc.expectedErr, err)
			}
		})
	}
}
