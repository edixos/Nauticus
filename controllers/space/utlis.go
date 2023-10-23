// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0
package space

import (
	"context"
	"errors"
	"reflect"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) FetchSpaceTemplate(ctx context.Context, name string) (*nauticusiov1alpha1.SpaceTemplate, error) {
	spaceTemplate := &nauticusiov1alpha1.SpaceTemplate{}

	err := r.Get(ctx, client.ObjectKey{
		Name: name,
	}, spaceTemplate)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// SpaceTemplate not found, return
			r.Log.Info("SpaceTemplate not found")
		}

		return nil, err
	}

	return spaceTemplate, nil
}

func (r *Reconciler) MergeResourceQuotas(space *nauticusiov1alpha1.Space, spacetpl *nauticusiov1alpha1.SpaceTemplate) (*corev1.ResourceQuotaSpec, error) {
	resourceQuotas := &corev1.ResourceQuotaSpec{}
	resourceQuotas.Hard = make(corev1.ResourceList)
	// Check if resourceQuota is provided in the Space and spaceTemplate
	switch {
	case !reflect.ValueOf(space.Spec.ResourceQuota).IsZero() && !reflect.ValueOf(spacetpl.Spec.ResourceQuota).IsZero():
		overrideResourceQuotas(resourceQuotas, space.Spec.ResourceQuota.Hard, spacetpl.Spec.ResourceQuota.Hard, corev1.ResourceLimitsCPU)
		overrideResourceQuotas(resourceQuotas, space.Spec.ResourceQuota.Hard, spacetpl.Spec.ResourceQuota.Hard, corev1.ResourceLimitsMemory)
		overrideResourceQuotas(resourceQuotas, space.Spec.ResourceQuota.Hard, spacetpl.Spec.ResourceQuota.Hard, corev1.ResourceRequestsCPU)
		overrideResourceQuotas(resourceQuotas, space.Spec.ResourceQuota.Hard, spacetpl.Spec.ResourceQuota.Hard, corev1.ResourceRequestsMemory)
	case reflect.ValueOf(space.Spec.ResourceQuota).IsZero() && !reflect.ValueOf(spacetpl.Spec.ResourceQuota).IsZero():
		resourceQuotas.Hard = spacetpl.Spec.ResourceQuota.Hard
	default:
		err := errors.New("both space and spacetpl resource quotas are empty. No merge required")

		return resourceQuotas, err
	}

	return resourceQuotas, nil
}

func (r *Reconciler) MergeRoleBindings(space *nauticusiov1alpha1.Space, spaceTemplate *nauticusiov1alpha1.SpaceTemplate) ([]nauticusiov1alpha1.AdditionalRoleBinding, error) {
	mergedRoleBindings := append([]nauticusiov1alpha1.AdditionalRoleBinding{}, space.Spec.AdditionalRoleBindings...)

	for _, roleBinding := range spaceTemplate.Spec.AdditionalRoleBindings {
		// Check if the role binding already exists in mergedRoleBindings
		if !cmpRoleBinding(mergedRoleBindings, roleBinding) {
			mergedRoleBindings = append(mergedRoleBindings, roleBinding)
		}
	}

	return mergedRoleBindings, nil
}

func (r *Reconciler) MergeNetworkPolicies(space *nauticusiov1alpha1.Space, spaceTemplate *nauticusiov1alpha1.SpaceTemplate) (nauticusiov1alpha1.NetworkPolicies, error) {
	mergedPolicies := space.Spec.NetworkPolicies

	// Check if EnableDefaultStrictMode is not provided in the space
	if reflect.ValueOf(mergedPolicies.EnableDefaultStrictMode).IsZero() {
		mergedPolicies.EnableDefaultStrictMode = spaceTemplate.Spec.NetworkPolicies.EnableDefaultStrictMode
	}

	for _, templatePolicy := range spaceTemplate.Spec.NetworkPolicies.Items {
		// Check if the policy already exists in mergedPolicies
		if !cmpNetworkPolicy(mergedPolicies, templatePolicy) {
			mergedPolicies.Items = append(mergedPolicies.Items, templatePolicy)
		}
	}

	return mergedPolicies, nil
}

func overrideResourceQuotas(resourceQuotas *corev1.ResourceQuotaSpec, spaceHard, templateHard corev1.ResourceList, resource corev1.ResourceName) {
	if spaceValue, exists := spaceHard[resource]; exists {
		resourceQuotas.Hard[resource] = spaceValue
	} else {
		resourceQuotas.Hard[resource] = templateHard[resource]
	}
}

func cmpRoleBinding(roleBindings []nauticusiov1alpha1.AdditionalRoleBinding, roleBinding nauticusiov1alpha1.AdditionalRoleBinding) bool {
	for _, rb := range roleBindings {
		if rb.RoleRef.Name == roleBinding.RoleRef.Name && rb.RoleRef.Kind == roleBinding.RoleRef.Kind {
			// Check if subjects are equal
			if reflect.DeepEqual(rb.Subjects, roleBinding.Subjects) {
				return true
			}
		}
	}

	return false
}

func cmpNetworkPolicy(mergedPolicies nauticusiov1alpha1.NetworkPolicies, policy networkingv1.NetworkPolicySpec) bool {
	for _, mergedPolicy := range mergedPolicies.Items {
		if reflect.DeepEqual(mergedPolicy, policy) {
			return true
		}
	}

	return false
}
