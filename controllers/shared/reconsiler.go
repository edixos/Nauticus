// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0
package shared

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciler reconciles an object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
}

// Clock is defined as a package var so it can be stubbed out during tests.
var Clock clock.Clock = clock.RealClock{}

// ConditionedObject is an interface that describes condition-related operations.
type ConditionedObject interface {
	client.Object
	metav1.Object
	GetConditions() []metav1.Condition
	SetConditions([]metav1.Condition)
}

func (r *Reconciler) EmitEvent(object runtime.Object, name string, res controllerutil.OperationResult, msg string, err error) {
	eventType := corev1.EventTypeNormal

	if err != nil {
		eventType = corev1.EventTypeWarning
		res = "Error"
	}

	r.Recorder.AnnotatedEventf(object, map[string]string{"OperationResult": string(res)}, eventType, name, msg)
}

func (r *Reconciler) setCondition(object ConditionedObject, observedGeneration int64, conditionType string, status metav1.ConditionStatus, reason, message string) {
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
	existingConditions := object.GetConditions()
	for idx, cond := range existingConditions {
		// Skip unrelated conditions
		if cond.Type != conditionType {
			continue
		}

		// If this update doesn't contain a state transition, we don't update
		// the conditions LastTransitionTime to Now()
		if cond.Status == status {
			newCondition.LastTransitionTime = cond.LastTransitionTime
		} else {
			r.Log.WithName(object.GetName()).Info("Found status change for Space condition, setting lastTransitionTime to", object.GetName(), nowTime)
		}

		// Overwrite the existing condition
		existingConditions[idx] = newCondition

		return
	}

	// If we've not found an existing condition of this type, we simply insert
	// the new condition into the slice.
	object.SetConditions(append(existingConditions, newCondition))
	r.Log.WithName(object.GetName()).Info("Setting lastTransitionTime condition for ", "name", object.GetName(), "time", newCondition.LastTransitionTime.Time)
}

func (r *Reconciler) ProcessCondition(ctx context.Context, object ConditionedObject, conditionType string, status metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, object.GetGeneration(), conditionType, status, reason, message)

	err := r.UpdateStatus(ctx, object)
	if err != nil {
		return
	}
}

func (r *Reconciler) UpdateStatus(ctx context.Context, object client.Object) (err error) {
	err = r.Client.Status().Update(ctx, object)
	if err != nil {
		r.Log.Info("Failed to update status", "object", object.GetName())

		return err
	}

	return nil
}
