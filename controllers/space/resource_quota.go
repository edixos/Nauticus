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

func (r *Reconciler) reconcileResourceQuota(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Status.NamespaceName,
		},
		Spec: space.Spec.ResourceQuota,
	}
	err = r.syncResourceQuotas(ctx, resourceQuota, space)

	return err
}

func (r *Reconciler) syncResourceQuotas(ctx context.Context, resourceQuota *corev1.ResourceQuota, space *nauticusiov1alpha1.Space) (err error) {
	var (
		res                            controllerutil.OperationResult
		spaceLabel, resourceQuotaLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if resourceQuotaLabel, err = v1alpha1.GetTypeLabel(resourceQuota); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, resourceQuota, func() error {
		resourceQuota.SetLabels(map[string]string{
			spaceLabel:         space.Name,
			resourceQuotaLabel: resourceQuota.Name,
		})
		resourceQuota.Spec = space.Spec.ResourceQuota

		return nil
	})
	r.Log.Info("ResourceQuota sync result: "+string(res), "name", resourceQuota.Name)
	r.EmitEvent(space, space.Name, res, "Ensuring ResourceQuota creation/Update", err)

	return err
}

func (r *Reconciler) deleteResourceQuota(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Status.NamespaceName,
		},
		Spec: space.Spec.ResourceQuota,
	}
	err = r.DeleteObject(ctx, resourceQuota)

	return err
}
