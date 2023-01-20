package controllers

import (
	"context"
	"reflect"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/go-logr/logr"
)

func (s *SpaceReconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
	log.Info("Reconciling Namespace for space.")
	err := s.reconcileNamespace(ctx, space, log)
	if err != nil {
		return err
	}
	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		log.Info("Reconciling Resource Quota for space")
		err = s.reconcileResourceQuota(ctx, space)
		if err != nil {
			return err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		log.Info("Reconciling Owner Role Binding for space")
		err = s.reconcileOwners(ctx, space, log)
		if err != nil {
			return err
		}
	}
	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalBindingSpecValue.IsZero() {
		log.Info("Reconciling Additional Role Binding for space")
		err = s.reconcileAdditionalRoleBindings(ctx, space, log)
		if err != nil {
			return err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPolicies.IsZero() {
		log.Info("Reconciling NetworkPolicies for space")
		err = s.reconcileNetworkPolicies(ctx, space)
		if err != nil {
			return err
		}
	}

	return nil
}
