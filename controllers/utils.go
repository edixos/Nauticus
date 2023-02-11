package controllers

import (
	"github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Clock is defined as a package var so it can be stubbed out during tests.
var Clock clock.Clock = clock.RealClock{}

func (s *SpaceReconciler) emitEvent(object runtime.Object, name string, res controllerutil.OperationResult, msg string, err error) {
	eventType := corev1.EventTypeNormal

	if err != nil {
		eventType = corev1.EventTypeWarning
		res = "Error"
	}

	s.Recorder.AnnotatedEventf(object, map[string]string{"OperationResult": string(res)}, eventType, name, msg)
}

func (s *SpaceReconciler) setSpaceCondition(space *v1alpha1.Space, observedGeneration int64, conditionType string, status metav1.ConditionStatus, reason, message string) {
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
	for idx, cond := range space.Status.Conditions {
		// Skip unrelated conditions
		if cond.Type != conditionType {
			continue
		}

		// If this update doesn't contain a state transition, we don't update
		// the conditions LastTransitionTime to Now()
		if cond.Status == status {
			newCondition.LastTransitionTime = cond.LastTransitionTime
		} else {
			s.Log.Info("Found status change for Space %q condition %q: %q -> %q; setting lastTransitionTime to %v", space.GetObjectMeta().GetName(), conditionType, cond.Status, status, nowTime.Time)
		}

		// Overwrite the existing condition
		space.Status.Conditions[idx] = newCondition

		return
	}

	// If we've not found an existing condition of this type, we simply insert
	// the new condition into the slice.
	space.Status.Conditions = append(space.Status.Conditions, newCondition)
	s.Log.Info("Setting lastTransitionTime for Space %q condition %q to %v", space.GetObjectMeta().GetName(), conditionType, nowTime.Time)
}
