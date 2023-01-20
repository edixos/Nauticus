package controllers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) emitEvent(object runtime.Object, name string, res controllerutil.OperationResult, msg string, err error) {
	eventType := corev1.EventTypeNormal

	if err != nil {
		eventType = corev1.EventTypeWarning
		res = "Error"
	}

	s.Recorder.AnnotatedEventf(object, map[string]string{"OperationResult": string(res)}, eventType, name, msg)
}
