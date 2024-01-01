// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package spacetemplate

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/controllers/shared"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Reconciler reconciles a SpaceTemplate object.
type Reconciler struct {
	shared.Reconciler
}

//+kubebuilder:rbac:groups=nauticus.io,resources=spacetemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nauticus.io,resources=spacetemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nauticus.io,resources=spacetemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the SpaceTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("spacetemplate", req.NamespacedName)

	// fetch the spaceTemplate
	spaceTpl := &nauticusiov1alpha1.SpaceTemplate{}

	err := r.Get(ctx, req.NamespacedName, spaceTpl)
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
		return r.reconcileDeleteSpaceTemplate(ctx, spaceTpl)
	}

	return r.reconcileSpaceTemplate(ctx, spaceTpl)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nauticusiov1alpha1.SpaceTemplate{}).
		Complete(r)
}
