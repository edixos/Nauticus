package controllers

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileServiceAccounts(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for _, serviceAccount := range space.Spec.ServiceAccounts.Items {
		sa := newServiceAccount(serviceAccount.Name, space.Status.NamespaceName, serviceAccount.Annotations)
		err = s.syncServiceAccount(ctx, sa, space, serviceAccount.Annotations)

		if err != nil {
			s.Log.Error(err, "Cannot Synchronize Service Account")

			return err
		}
	}

	return nil
}

func (s *SpaceReconciler) syncServiceAccount(ctx context.Context, serviceAccount *corev1.ServiceAccount, space *nauticusiov1alpha1.Space, annotations nauticusiov1alpha1.Annotations) (err error) {
	var (
		res                             controllerutil.OperationResult
		spaceLabel, serviceAccountLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if serviceAccountLabel, err = v1alpha1.GetTypeLabel(serviceAccount); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, s.Client, serviceAccount, func() (err error) {
		serviceAccount.SetLabels(map[string]string{
			spaceLabel:          space.Name,
			serviceAccountLabel: serviceAccount.Name,
		})
		serviceAccount.SetAnnotations(annotations)

		return controllerutil.SetControllerReference(space, serviceAccount, s.Client.Scheme())
	})
	s.Log.Info("ServiceAccount sync result: "+string(res), "name", serviceAccount.Name, "namespace", space.Status.NamespaceName)
	s.emitEvent(space, space.Name, res, "Ensuring ServiceAccount creation/Update", err)

	return nil
}

func newServiceAccount(name, namespace string, annotations nauticusiov1alpha1.Annotations) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
	}
}
