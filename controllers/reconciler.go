// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"reflect"
	"time"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NauticusSpaceFinalizer = "nauticus.io/finalizer"
	requeueAfter           = time.Minute * 3

	SpaceConditionReady    nauticusiov1alpha1.ConditionType = "Ready"
	SpaceConditionCreating nauticusiov1alpha1.ConditionType = "Creating"
	SpaceConditionFailed   nauticusiov1alpha1.ConditionType = "Failed"

	SpaceConditionStatusUnknown = metav1.ConditionStatus(corev1.ConditionUnknown)
	SpaceConditionStatusTrue    = metav1.ConditionStatus(corev1.ConditionTrue)
	SpaceConditionStatusFalse   = metav1.ConditionStatus(corev1.ConditionFalse)

	SpaceSyncSuccessReason nauticusiov1alpha1.ConditionReason = "SpaceSyncedSuccessfully"
	SpaceCreatingReason    nauticusiov1alpha1.ConditionReason = "SpaceCreating"
	SpaceFailedReason      nauticusiov1alpha1.ConditionReason = "SpaceSyncFailed"

	SpaceSyncSuccessMessage nauticusiov1alpha1.ConditionMessage = "Space synced successfully"
	SpaceSyncFailMessage    nauticusiov1alpha1.ConditionMessage = "Space failed to sync"
	SpaceCreatingMessage    nauticusiov1alpha1.ConditionMessage = "Creating Space in progress"
)

func (s *SpaceReconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	if !controllerutil.ContainsFinalizer(space, NauticusSpaceFinalizer) {
		controllerutil.AddFinalizer(space, NauticusSpaceFinalizer)

		if err = s.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}

	s.Log.Info("Reconciling Namespace for space.")

	s.processInProgressCondition(ctx, space)
	s.setMetrics(space, SpaceConditionCreating)

	err = s.reconcileNamespace(ctx, space)
	if err != nil {
		s.processFailedCondition(ctx, space)
		s.setMetrics(space, SpaceConditionFailed)

		return ctrl.Result{}, err
	}
	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		s.Log.Info("Reconciling Resource Quota for space")
		err = s.reconcileResourceQuota(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Owner Role Binding for space")
		err = s.reconcileOwners(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Additional Role Binding for space")
		err = s.reconcileAdditionalRoleBindings(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPolicies.IsZero() {
		s.Log.Info("Reconciling NetworkPolicies for space")
		err = s.reconcileNetworkPolicies(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	limitRanges := reflect.ValueOf(space.Spec.LimitRanges)
	if !limitRanges.IsZero() {
		s.Log.Info("Reconciling LimitRanges for space")
		err = s.reconcileLimitRanges(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	serviceAccounts := reflect.ValueOf(space.Spec.ServiceAccounts)
	if !serviceAccounts.IsZero() {
		s.Log.Info("Reconciling ServiceAccounts for space")
		err = s.reconcileServiceAccounts(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)
			s.setMetrics(space, SpaceConditionFailed)

			return ctrl.Result{}, err
		}
	}

	s.processReadyCondition(ctx, space)
	s.setMetrics(space, SpaceConditionReady)

	return ctrl.Result{
		RequeueAfter: requeueAfter,
	}, nil
}

func (s *SpaceReconciler) reconcileSpaceFromTemplate(ctx context.Context, req ctrl.Request, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {

	// Fetch data from the SpaceTemplate
	spacetpl, err := s.FetchSpaceTemplate(ctx, req, space.Spec.Template.Name, space.Spec.Template.Namespace)

	// Update the existing Space resource with the data from the SpaceTemplate
	// Check if specific fields in the Space spec are not provided
	// todo : proper override(compare & append)
	if reflect.ValueOf(space.Spec.ResourceQuota).IsZero() {
		// use  the value provided in the SpaceTemplate
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
	s.Log.Info("Reconciling Space from", "SpaceTemplate", spacetpl.Name)

	return s.reconcileSpace(ctx, space)
}

func (s *SpaceReconciler) reconcileDelete(ctx context.Context, space *nauticusiov1alpha1.Space) (result reconcile.Result, err error) {
	// The annotation is set, so skip namespace deletion
	// Just remove the finalizer from the Space
	if space.HasIgnoreUnderlyingDeletionAnnotation() {
		if controllerutil.ContainsFinalizer(space, NauticusSpaceFinalizer) {
			controllerutil.RemoveFinalizer(space, NauticusSpaceFinalizer)

			if err = s.Update(ctx, space); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, err
	}
	// If the annotation is not set, delete all created resources
	if controllerutil.ContainsFinalizer(space, NauticusSpaceFinalizer) {
		if err = s.deleteNetworkPolicies(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteLimitRanges(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteOwners(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteAdditionalRoleBindings(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteResourceQuota(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteServiceAccounts(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = s.deleteNamespace(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		// remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(space, NauticusSpaceFinalizer)

		if err = s.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Stop reconciliation as the item is being deleted
	return ctrl.Result{}, err
}
