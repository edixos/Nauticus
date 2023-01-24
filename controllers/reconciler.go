package controllers

import (
	"context"
	"reflect"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NauticusFinalizer = "nauticus.io/finalizer"
)

func (s *SpaceReconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	if !controllerutil.ContainsFinalizer(space, NauticusFinalizer) {
		controllerutil.AddFinalizer(space, NauticusFinalizer)

		if err = s.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}

	s.Log.Info("Reconciling Namespace for space.")

	err = s.reconcileNamespace(ctx, space)
	if err != nil {
		return ctrl.Result{}, err
	}

	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		s.Log.Info("Reconciling Resource Quota for space")
		err = s.reconcileResourceQuota(ctx, space)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Owner Role Binding for space")
		err = s.reconcileOwners(ctx, space)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Additional Role Binding for space")
		err = s.reconcileAdditionalRoleBindings(ctx, space)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPolicies.IsZero() {
		s.Log.Info("Reconciling NetworkPolicies for space")
		err = s.reconcileNetworkPolicies(ctx, space)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (s *SpaceReconciler) reconcileDelete(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	if controllerutil.ContainsFinalizer(space, NauticusFinalizer) {
		namespace, _ := s.newNamespace(space)
		err = s.Client.Delete(ctx, namespace)

		// remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(space, NauticusFinalizer)

		if err = s.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Stop reconciliation as the item is being deleted
	return ctrl.Result{}, nil
}
