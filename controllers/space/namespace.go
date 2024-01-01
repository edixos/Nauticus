// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package space

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) reconcileNamespace(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	namespace := r.newNamespace(space)

	err = r.syncNamespace(ctx, namespace, space)

	if err != nil {
		return err
	}

	// Update the Space's status
	space.Status.NamespaceName = namespace.Name
	err = r.UpdateStatus(ctx, space)

	if err != nil {
		r.Log.Info("error updating the status")
	}

	return err
}

func (r *Reconciler) syncNamespace(ctx context.Context, namespace *corev1.Namespace, space *nauticusiov1alpha1.Space) (err error) {
	var (
		res                        controllerutil.OperationResult
		spaceLabel, namespaceLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if namespaceLabel, err = v1alpha1.GetTypeLabel(namespace); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, namespace, func() error {
		namespace.SetLabels(map[string]string{
			spaceLabel:     space.Name,
			namespaceLabel: namespace.Name,
		})

		return nil
	})
	r.Log.Info("Namespace sync result: "+string(res), "name", namespace.Name)
	r.EmitEvent(space, space.Name, res, "Ensuring Namespace creation", err)

	return err
}

func (r *Reconciler) newNamespace(space *nauticusiov1alpha1.Space) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   r.namespaceName(space),
			Labels: space.Labels,
		},
	}

	return namespace
}

func (r *Reconciler) namespaceName(space *nauticusiov1alpha1.Space) string {
	return space.Name
}

func (r *Reconciler) deleteNamespace(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	namespace := r.newNamespace(space)

	return r.DeleteObject(ctx, namespace)
}
