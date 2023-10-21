// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package spacetemplate

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/controller/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *Reconciler) reconcileSpaceTemplate(ctx context.Context, spaceTpl *nauticusiov1alpha1.SpaceTemplate) (result reconcile.Result, err error) {
	if !controllerutil.ContainsFinalizer(spaceTpl, constants.NauticusSpaceFinalizer) {
		controllerutil.AddFinalizer(spaceTpl, constants.NauticusSpaceFinalizer)

		if err = r.Update(ctx, spaceTpl); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.ProcessCondition(
		ctx,
		spaceTpl,
		constants.SpaceTplConditionReady,
		metav1.ConditionTrue,
		constants.SpaceTplSyncSuccessReason,
		constants.SpaceTplSyncSuccessMessage,
	)
	r.EmitEvent(spaceTpl, spaceTpl.GetName(), controllerutil.OperationResultCreated, constants.SpaceTplCreatingMessage, nil)

	return ctrl.Result{
		RequeueAfter: constants.RequeueAfter,
	}, nil
}

func (r *Reconciler) reconcileDeleteSpaceTemplate(ctx context.Context, spaceTpl *nauticusiov1alpha1.SpaceTemplate) (result reconcile.Result, err error) {
	if controllerutil.ContainsFinalizer(spaceTpl, constants.NauticusSpaceFinalizer) {
		controllerutil.RemoveFinalizer(spaceTpl, constants.NauticusSpaceFinalizer)

		if err = r.Update(ctx, spaceTpl); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, err
}
