package controllers

import (
	"context"

	"github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/metrics"
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
			s.Log.WithName(space.GetName()).Info("Found status change for Space condition, setting lastTransitionTime to", space.GetName(), nowTime)
		}

		// Overwrite the existing condition
		space.Status.Conditions[idx] = newCondition

		return
	}

	// If we've not found an existing condition of this type, we simply insert
	// the new condition into the slice.
	space.Status.Conditions = append(space.Status.Conditions, newCondition)
	s.Log.WithName(space.GetName()).Info("Setting lastTransitionTime for Space condition ", space.GetObjectMeta().GetName(), nowTime.Time)
}

func (s *SpaceReconciler) processFailedCondition(ctx context.Context, space *v1alpha1.Space) {
	s.setSpaceCondition(
		space,
		space.GetGeneration(),
		string(SpaceConditionFailed),
		SpaceConditionStatusFalse,
		string(SpaceFailedReason),
		string(SpaceSyncFailMessage),
	)
	s.updateStatus(ctx, space)
}

func (s *SpaceReconciler) processReadyCondition(ctx context.Context, space *v1alpha1.Space) {
	s.setSpaceCondition(
		space,
		space.GetGeneration(),
		string(SpaceConditionReady),
		SpaceConditionStatusTrue,
		string(SpaceSyncSuccessReason),
		string(SpaceSyncSuccessMessage),
	)
	s.updateStatus(ctx, space)
}

func (s *SpaceReconciler) processInProgressCondition(ctx context.Context, space *v1alpha1.Space) {
	s.setSpaceCondition(
		space,
		space.GetGeneration(),
		string(SpaceConditionCreating),
		SpaceConditionStatusUnknown,
		string(SpaceCreatingReason),
		string(SpaceCreatingMessage),
	)
	s.updateStatus(ctx, space)
}

func (s *SpaceReconciler) updateStatus(ctx context.Context, space *v1alpha1.Space) {
	err := s.Client.Status().Update(ctx, space)
	if err != nil {
		s.Log.Info("Failed to update Space status", "space", space.Name)
	}
}

func (s *SpaceReconciler) setMetrics(space *v1alpha1.Space, conditionType v1alpha1.ConditionType) {
	switch conditionType {
	case SpaceConditionCreating:
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(0)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(1)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(0)
	case SpaceConditionReady:
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(1)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(0)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(0)
	case SpaceConditionFailed:
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(0)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(0)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(1)
	}
}
