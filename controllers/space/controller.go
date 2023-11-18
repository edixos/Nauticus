// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package space

import (
	"context"

	"github.com/edixos/nauticus/controllers/shared"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Reconciler reconciles a Space object.
type Reconciler struct {
	shared.Reconciler
}

//+kubebuilder:rbac:groups=nauticus.io,resources=spaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nauticus.io,resources=spaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nauticus.io,resources=spaces/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=limitranges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=bind

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Space object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("space", req.NamespacedName)

	// Fetch the Space instance
	space := &nauticusiov1alpha1.Space{}

	err := r.Get(ctx, req.NamespacedName, space)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Space not found, return
			log.Info("Space not found.")

			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if !space.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, space)
	}

	if space.Spec.TemplateRef.Name != "" {
		return r.reconcileSpaceFromTemplate(ctx, space)
	}

	return r.reconcileSpace(ctx, space)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nauticusiov1alpha1.Space{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}

func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}
