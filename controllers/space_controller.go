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

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SpaceReconciler reconciles a Space object
type SpaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=nauticus.io,resources=spaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nauticus.io,resources=spaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nauticus.io,resources=spaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Space object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (s *SpaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := s.Log.WithValues("space", req.NamespacedName)
	ctx = context.TODO()

	// Fetch the Space instance
	space := &nauticusiov1alpha1.Space{}
	err := s.Get(ctx, req.NamespacedName, space)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Space not found, return
			log.Info("Space not found.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if err := s.reconcileSpace(ctx, space, log); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (s *SpaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nauticusiov1alpha1.Space{}).
		Owns(&v1.Namespace{}).
		Owns(&v1.ResourceQuota{}).
		Complete(s)
}
