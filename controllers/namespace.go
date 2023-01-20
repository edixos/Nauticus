package controllers

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileNamespace(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
	namespace, err := s.newNamespace(space)
	if err != nil {
		return err
	}
	err = s.syncNamespace(ctx, namespace, space)
	// Update the Space's status
	space.Status.NamespaceName = namespace.Name
	err = s.Client.Status().Update(ctx, space)
	if err != nil {
		log.Error(err, "Failed to update Space status", "space", space.Name)
		return err
	}
	return nil
}

func (s *SpaceReconciler) syncNamespace(ctx context.Context, namespace *corev1.Namespace, space *nauticusiov1alpha1.Space) (err error) {
	var res controllerutil.OperationResult
	var spaceLabel, namespaceLabel string
	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}
	if namespaceLabel, err = v1alpha1.GetTypeLabel(namespace); err != nil {
		return
	}
	res, err = controllerutil.CreateOrUpdate(ctx, s.Client, namespace, func() error {
		namespace.SetLabels(map[string]string{
			spaceLabel:     space.Name,
			namespaceLabel: namespace.Name,
		})
		return controllerutil.SetControllerReference(space, namespace, s.Scheme)
	})
	s.Log.Info("Namespace sync result: "+string(res), "name", namespace.Name)
	return nil
}

func (s *SpaceReconciler) newNamespace(space *nauticusiov1alpha1.Space) (*corev1.Namespace, error) {

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.namespaceName(space),
			Labels: space.Labels,
		},
	}
	return namespace, nil
}

func (s *SpaceReconciler) namespaceName(space *nauticusiov1alpha1.Space) string {
	return space.Name
}
