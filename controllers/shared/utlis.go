// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) DeleteObject(ctx context.Context, object client.Object) (err error) {
	if err = r.Client.Delete(ctx, object); client.IgnoreNotFound(err) != nil {
		return err
	}

	return nil
}

func (r *Reconciler) ProcessFailedCondition(ctx context.Context, object ConditionedObject, conditionType string, status metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, object.GetGeneration(), conditionType, status, reason, message)
	err := r.UpdateStatus(ctx, object)
	if err != nil { //nolint
		return
	}
}

func (r *Reconciler) ProcessReadyCondition(ctx context.Context, object ConditionedObject, conditionType string, status metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, object.GetGeneration(), conditionType, status, reason, message)
	err := r.UpdateStatus(ctx, object)
	if err != nil { //nolint
		return
	}
}

func (r *Reconciler) ProcessInProgressCondition(ctx context.Context, object ConditionedObject, conditionType string, status metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, object.GetGeneration(), conditionType, status, reason, message)
	err := r.UpdateStatus(ctx, object)
	if err != nil { //nolint
		return
	}
}
