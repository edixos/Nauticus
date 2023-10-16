/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SpaceTplConditionReady nauticusiov1alpha1.ConditionType = "Ready"
	//	SpaceTplConditionCreating nauticusiov1alpha1.ConditionType = "Creating"
	//	SpaceTplConditionFailed   nauticusiov1alpha1.ConditionType = "Failed"

	//	SpaceTplConditionStatusUnknown = metav1.ConditionStatus(corev1.ConditionUnknown)
	SpaceTplConditionStatusTrue = metav1.ConditionStatus(corev1.ConditionTrue)
	//	SpaceTplConditionStatusFalse   = metav1.ConditionStatus(corev1.ConditionFalse)

	//	SpaceTplSyncSuccessReason nauticusiov1alpha1.ConditionReason = "SpaceTemplateSyncedSuccessfully"
	SpaceTplCreatingReason nauticusiov1alpha1.ConditionReason = "SpaceTemplateCreating"
	//SpaceTplFailedReason      nauticusiov1alpha1.ConditionReason = "SpaceTemplateSyncFailed"

	//	SpaceTplSyncSuccessMessage nauticusiov1alpha1.ConditionMessage = "SpaceTemplate synced successfully"
	//	SpaceTplSyncFailMessage    nauticusiov1alpha1.ConditionMessage = "SpaceTemplate failed to sync"
	SpaceTplCreatingMessage nauticusiov1alpha1.ConditionMessage = "Creating SpaceTemplate in progress"
)

// SpaceTemplateReconciler reconciles a SpaceTemplate object
type SpaceTemplateReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
}

//+kubebuilder:rbac:groups=nauticus.io.nauticus.io,resources=spacetemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nauticus.io.nauticus.io,resources=spacetemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nauticus.io.nauticus.io,resources=spacetemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the SpaceTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (st *SpaceTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := st.Log.WithValues("spacetemplate", req.NamespacedName)

	// fetch the spaceTemplate
	spaceTpl := &nauticusiov1alpha1.SpaceTemplate{}

	err := st.Get(ctx, req.NamespacedName, spaceTpl)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Space not found, return
			log.Info("SpaceTemplate not found.")

			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}
	if !spaceTpl.ObjectMeta.DeletionTimestamp.IsZero() {
		return st.reconcileDeleteSpaceTemplate(ctx, spaceTpl)
	}

	return st.reconcileSpaceTemplate(ctx, spaceTpl)
}

// SetupWithManager sets up the controller with the Manager.
func (st *SpaceTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nauticusiov1alpha1.SpaceTemplate{}).
		Complete(st)
}

func (st *SpaceTemplateReconciler) reconcileDeleteSpaceTemplate(ctx context.Context, spaceTpl *nauticusiov1alpha1.SpaceTemplate) (result reconcile.Result, err error) {
	if controllerutil.ContainsFinalizer(spaceTpl, NauticusSpaceFinalizer) {
		controllerutil.RemoveFinalizer(spaceTpl, NauticusSpaceFinalizer)
		if err = st.Update(ctx, spaceTpl); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, err
}

func (st *SpaceTemplateReconciler) reconcileSpaceTemplate(ctx context.Context, spaceTpl *nauticusiov1alpha1.SpaceTemplate) (result reconcile.Result, err error) {
	if !controllerutil.ContainsFinalizer(spaceTpl, NauticusSpaceFinalizer) {
		controllerutil.AddFinalizer(spaceTpl, NauticusSpaceFinalizer)

		if err = st.Update(ctx, spaceTpl); err != nil {
			return ctrl.Result{}, err
		}
	}
	st.Log.Info("Reconciling spaceTemplate")

	return ctrl.Result{
		RequeueAfter: requeueAfter,
	}, nil
}

func (st *SpaceTemplateReconciler) processInProgressCondition(ctx context.Context, spaceTpl *nauticusiov1alpha1.SpaceTemplate) {
	st.setSpaceConditionSpaceTpl(
		spaceTpl,
		spaceTpl.GetGeneration(),
		string(SpaceTplConditionReady),
		SpaceTplConditionStatusTrue,
		string(SpaceTplCreatingReason),
		string(SpaceTplCreatingMessage),
	)
	st.updateStatusSpaceTpl(ctx, spaceTpl)
}

func (st *SpaceTemplateReconciler) setSpaceConditionSpaceTpl(spaceTpl *nauticusiov1alpha1.SpaceTemplate, observedGeneration int64, conditionType string, status metav1.ConditionStatus, reason, message string) {
	newCondition := metav1.Condition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	}
	nowTime := metav1.NewTime(Clock.Now())
	newCondition.LastTransitionTime = nowTime

	// Set the condition generation
	newCondition.ObservedGeneration = observedGeneration

	// Search through existing conditions
	for idx, cond := range spaceTpl.Status.Conditions {
		// Skip unrelated conditions
		if cond.Type != conditionType {
			continue
		}

		// If this update doesn't contain a state transition, we don't update
		// the conditions LastTransitionTime to Now()
		if cond.Status == status {
			newCondition.LastTransitionTime = cond.LastTransitionTime
		} else {
			st.Log.WithName(spaceTpl.GetName()).Info("Found status change for Space condition, setting lastTransitionTime to", spaceTpl.GetName(), nowTime)
		}

		// Overwrite the existing condition
		spaceTpl.Status.Conditions[idx] = newCondition

		return
	}

	// If we've not found an existing condition of this type, we simply insert
	// the new condition into the slice.
	spaceTpl.Status.Conditions = append(spaceTpl.Status.Conditions, newCondition)
	st.Log.WithName(spaceTpl.GetName()).Info("Setting lastTransitionTime for Space condition ", spaceTpl.GetObjectMeta().GetName(), nowTime.Time)
}

func (st *SpaceTemplateReconciler) updateStatusSpaceTpl(ctx context.Context, spacetpl *nauticusiov1alpha1.SpaceTemplate) {
	err := st.Client.Status().Update(ctx, spacetpl)
	if err != nil {
		st.Log.Info("Failed to update SpaceTemplate status", "SpaceTemplate", spacetpl.Name)
	}
}
