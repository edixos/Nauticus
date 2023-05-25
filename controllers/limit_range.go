package controllers

import (
	"context"
	"strconv"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileLimitRanges(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for i, limitRange := range space.Spec.LimitRanges.Items {
		lrName := "nauticus-custom-" + strconv.Itoa(i)
		lr := newLimitRange(lrName, space.Status.NamespaceName, limitRange)
		err = s.syncLimitRange(ctx, lr, space, limitRange)

		if err != nil {
			s.Log.Error(err, "Cannot Synchronize Limit Range")

			return err
		}
	}

	return nil
}

func (s *SpaceReconciler) syncLimitRange(ctx context.Context, limitRange *corev1.LimitRange, space *nauticusiov1alpha1.Space, spec corev1.LimitRangeSpec) (err error) {
	var (
		res                         controllerutil.OperationResult
		spaceLabel, limitRangeLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if limitRangeLabel, err = v1alpha1.GetTypeLabel(limitRange); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, s.Client, limitRange, func() (err error) {
		limitRange.SetLabels(map[string]string{
			spaceLabel:      space.Name,
			limitRangeLabel: limitRange.Name,
		})
		limitRange.Spec = spec

		return nil
	})
	s.Log.Info("LimitRange sync result: "+string(res), "name", limitRange.Name, "namespace", space.Status.NamespaceName)
	s.emitEvent(space, space.Name, res, "Ensuring LimitRange creation/Update", err)

	return nil
}

func newLimitRange(name, namespace string, limitRangeSpec corev1.LimitRangeSpec) *corev1.LimitRange {
	return &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: limitRangeSpec,
	}
}

func (s *SpaceReconciler) deleteLimitRanges(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for i, limitRange := range space.Spec.LimitRanges.Items {
		lrName := "nauticus-custom-" + strconv.Itoa(i)
		lr := newLimitRange(lrName, space.Status.NamespaceName, limitRange)

		if err = s.deleteObject(ctx, lr); err != nil {
			return err
		}
	}

	return nil
}
