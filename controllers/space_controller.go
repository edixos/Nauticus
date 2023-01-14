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
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/util/rand"
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Space object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *SpaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("space", req.NamespacedName)
    ctx = context.Background()

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
    // Generate random suffix
    suffix := rand.String(4)
    // Create a new namespace object
    namespace := &v1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name:   space.Name + "-" + suffix,
            Labels: space.Labels,
        },
    }

    if err := ctrl.SetControllerReference(space, namespace, r.Scheme); err != nil {
        return ctrl.Result{}, err
    }

    // Check if the namespace exist and create it if it does not
    existingNamespace := &v1.Namespace{}
    err = r.Client.Get(ctx, client.ObjectKey{Name: space.Status.NamespaceName}, existingNamespace)
    if err != nil {
        if apierrors.IsNotFound(err) {
            // Namespace does not exist, creating it.
            log.Info("Creating the namespace", "namespace", namespace.Name)
            err = r.Client.Create(ctx, namespace)
            if err != nil {
                log.Error(err, "Failed to create namespace", "namespace", namespace.Name)
                return ctrl.Result{}, err
            }
            log.Info("The namespace created successfully.", "namespace", namespace.Name)
        } else {
            log.Error(err, "Failed to check if the namespace exists", "namespace", namespace.Name)
        }
    } else {
        namespace = existingNamespace
    }

    // Update the Space's status
    space.Status.NamespaceName = namespace.Name
    err = r.Client.Status().Update(ctx, space)
    if err != nil {
        log.Error(err, "Failed to update Space status", "space", space.Name)
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&nauticusiov1alpha1.Space{}).
        Owns(&v1.Namespace{}).
        Complete(r)
}
