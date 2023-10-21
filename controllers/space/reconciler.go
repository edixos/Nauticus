// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package space

import (
	"context"
	"reflect"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/controller/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *Reconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	if !controllerutil.ContainsFinalizer(space, constants.NauticusSpaceFinalizer) {
		controllerutil.AddFinalizer(space, constants.NauticusSpaceFinalizer)

		if err = r.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Log.Info("Reconciling Namespace for space.")

	r.ProcessInProgressCondition(ctx, space, constants.SpaceConditionCreating, metav1.ConditionUnknown, constants.SpaceCreatingReason, constants.SpaceCreatingMessage)
	r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionCreating))

	err = r.reconcileNamespace(ctx, space)
	if err != nil {
		r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
		r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

		return ctrl.Result{}, err
	}

	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		r.Log.Info("Reconciling Resource Quota for space")
		err = r.reconcileResourceQuota(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		r.Log.Info("Reconciling Owner Role Binding for space")
		err = r.reconcileOwners(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalBindingSpecValue.IsZero() {
		r.Log.Info("Reconciling Additional Role Binding for space")
		err = r.reconcileAdditionalRoleBindings(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPolicies.IsZero() {
		r.Log.Info("Reconciling NetworkPolicies for space")
		err = r.reconcileNetworkPolicies(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	limitRanges := reflect.ValueOf(space.Spec.LimitRanges)
	if !limitRanges.IsZero() {
		r.Log.Info("Reconciling LimitRanges for space")
		err = r.reconcileLimitRanges(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	serviceAccounts := reflect.ValueOf(space.Spec.ServiceAccounts)
	if !serviceAccounts.IsZero() {
		r.Log.Info("Reconciling ServiceAccounts for space")
		err = r.reconcileServiceAccounts(ctx, space)

		if err != nil {
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionFailed))

			return ctrl.Result{}, err
		}
	}

	r.ProcessReadyCondition(ctx, space, constants.SpaceConditionReady, metav1.ConditionTrue, constants.SpaceSyncSuccessReason, constants.SpaceSyncSuccessMessage)
	r.setMetrics(space, nauticusiov1alpha1.ConditionType(constants.SpaceConditionReady))

	return ctrl.Result{
		RequeueAfter: constants.RequeueAfter,
	}, nil
}

func (r *Reconciler) reconcileSpaceFromTemplate(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	// Fetch data from the SpaceTemplate
	spacetpl, err := r.FetchSpaceTemplate(ctx, space.Spec.TemplateRef.Name)
	if err != nil {
		r.Log.Info("SpaceTemplate:", "SpaceTemplate does not exist", space.Spec.TemplateRef.Name)

		return ctrl.Result{}, err
	}

	// Update the existing Space resource with the data from the SpaceTemplate
	// Check if specific fields in the Space spec are not provided
	if reflect.ValueOf(space.Spec.ResourceQuota).IsZero() {
		space.Spec.ResourceQuota = spacetpl.Spec.ResourceQuota
	}

	if reflect.ValueOf(space.Spec.AdditionalRoleBindings).IsZero() {
		space.Spec.AdditionalRoleBindings = spacetpl.Spec.AdditionalRoleBindings
	}

	if reflect.ValueOf(space.Spec.NetworkPolicies).IsZero() {
		space.Spec.NetworkPolicies = spacetpl.Spec.NetworkPolicies
	}

	if reflect.ValueOf(space.Spec.LimitRanges).IsZero() {
		space.Spec.LimitRanges = spacetpl.Spec.LimitRanges
	}
	// Create or update the Space in the cluster
	r.Log.Info("Reconciling Space from", "SpaceTemplate", spacetpl.Name)

	return r.reconcileSpace(ctx, space)
}

func (r *Reconciler) reconcileDelete(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	// The annotation is set, so skip namespace deletion
	// Just remove the finalizer from the Space
	if space.HasIgnoreUnderlyingDeletionAnnotation() {
		if controllerutil.ContainsFinalizer(space, constants.NauticusSpaceFinalizer) {
			controllerutil.RemoveFinalizer(space, constants.NauticusSpaceFinalizer)

			if err = r.Update(ctx, space); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, err
	}
	// If the annotation is not set, delete all created resources
	if controllerutil.ContainsFinalizer(space, constants.NauticusSpaceFinalizer) {
		if err = r.deleteNetworkPolicies(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteLimitRanges(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteOwners(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteAdditionalRoleBindings(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteResourceQuota(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteServiceAccounts(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteNamespace(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		// remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(space, constants.NauticusSpaceFinalizer)

		if err = r.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Stop reconciliation as the item is being deleted
	return ctrl.Result{}, err
}
