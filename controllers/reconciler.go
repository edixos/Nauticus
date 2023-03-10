package controllers

import (
	"context"
	"reflect"
	"time"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NauticusFinalizer = "nauticus.io/finalizer"
	requeueAfter      = time.Minute * 3

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
	if !controllerutil.ContainsFinalizer(space, NauticusFinalizer) {
		controllerutil.AddFinalizer(space, NauticusFinalizer)

		if err = s.Update(ctx, space); err != nil {
			return ctrl.Result{}, err
		}
	}

	s.Log.Info("Reconciling Namespace for space.")

	s.processInProgressCondition(ctx, space)

	err = s.reconcileNamespace(ctx, space)
	if err != nil {
		s.processFailedCondition(ctx, space)

		return ctrl.Result{}, err
	}

	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		s.Log.Info("Reconciling Resource Quota for space")
		err = s.reconcileResourceQuota(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Owner Role Binding for space")
		err = s.reconcileOwners(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	additionalBindingSpecValue := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalBindingSpecValue.IsZero() {
		s.Log.Info("Reconciling Additional Role Binding for space")
		err = s.reconcileAdditionalRoleBindings(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	networkPolicies := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPolicies.IsZero() {
		s.Log.Info("Reconciling NetworkPolicies for space")
		err = s.reconcileNetworkPolicies(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	limitRanges := reflect.ValueOf(space.Spec.LimitRanges)
	if !limitRanges.IsZero() {
		s.Log.Info("Reconciling LimitRanges for space")
		err = s.reconcileLimitRanges(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	serviceAccounts := reflect.ValueOf(space.Spec.ServiceAccounts)
	if !serviceAccounts.IsZero() {
		s.Log.Info("Reconciling ServiceAccounts for space")
		err = s.reconcileServiceAccounts(ctx, space)

		if err != nil {
			s.processFailedCondition(ctx, space)

			return ctrl.Result{}, err
		}
	}

	s.processReadyCondition(ctx, space)

	return ctrl.Result{
		RequeueAfter: requeueAfter,
	}, nil
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
