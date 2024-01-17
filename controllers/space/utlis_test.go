// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0
package space

import (
	"reflect"
	"testing"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestMergeResourceQuotas(t *testing.T) {
	testCases := map[string]struct {
		Name          string
		Space         *nauticusiov1alpha1.Space
		SpaceTemplate *nauticusiov1alpha1.SpaceTemplate
		Expected      *corev1.ResourceQuotaSpec
		ExpectedErr   error
	}{
		"BothResourceQuotasProvided": {
			Name: "Both space and spaceTemplate ResourceQuotas provided",
			Space: &nauticusiov1alpha1.Space{
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
			SpaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
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
			Expected: &corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceLimitsCPU:      resource.MustParse("8"),
					corev1.ResourceLimitsMemory:   resource.MustParse("16Gi"),
					corev1.ResourceRequestsCPU:    resource.MustParse("4"),
					corev1.ResourceRequestsMemory: resource.MustParse("8Gi"),
				},
			},
			ExpectedErr: nil,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := MergeResourceQuotas(tc.Space, tc.SpaceTemplate)

			if !reflect.DeepEqual(tc.Expected, result) {
				t.Errorf("Expected: %v, Got: %v", tc.Expected, result)
			}
			if !reflect.DeepEqual(tc.ExpectedErr, err) {
				t.Errorf("Expected error: %v, Got error: %v", tc.ExpectedErr, err)
			}
		})
	}
}

func TestMergeRoleBindings(t *testing.T) {
	testCases := map[string]struct {
		Name          string
		Space         *nauticusiov1alpha1.Space
		SpaceTemplate *nauticusiov1alpha1.SpaceTemplate
		Expected      []nauticusiov1alpha1.AdditionalRoleBinding
		ExpectedErr   error
	}{
		"BothRoleBindingsProvided": {
			Name: "Space has role bindings, SpaceTemplate has role bindings",
			Space: &nauticusiov1alpha1.Space{
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
			SpaceTemplate: &nauticusiov1alpha1.SpaceTemplate{
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
			Expected: []nauticusiov1alpha1.AdditionalRoleBinding{
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
			ExpectedErr: nil,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := MergeRoleBindings(tc.Space, tc.SpaceTemplate)

			if !reflect.DeepEqual(tc.Expected, result) {
				t.Errorf("Expected: %v, Got: %v", tc.Expected, result)
			}
			if !reflect.DeepEqual(tc.ExpectedErr, err) {
				t.Errorf("Expected error: %v, Got error: %v", tc.ExpectedErr, err)
			}
		})
	}
}
