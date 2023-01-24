package controllers

import (
	"context"
	"reflect"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
)

func (s *SpaceReconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space) error {
	s.Log.Info("Reconciling Namespace for space.")

	err := s.reconcileNamespace(ctx, space)
	if err != nil {
		return err
	}

	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		s.Log.Info("Reconciling Resource Quota for space")
		err = s.reconcileResourceQuota(ctx, space)

		if err != nil {
			return err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)

	if !ownerRoleBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Owner Role Binding for space")
		err = s.reconcileOwners(ctx, space)

		if err != nil {
			return err
		}
	}

	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)

	if !additionalBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Additional Role Binding for space")
		err = s.reconcileAdditionalRoleBindings(ctx, space)

		if err != nil {
			return err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)

	if !networkPolicies.IsZero() {
		s.Log.Info("Reconciling NetworkPolicies for space")
		err = s.reconcileNetworkPolicies(ctx, space)

		if err != nil {
			return err
		}
	}

	return nil
}
