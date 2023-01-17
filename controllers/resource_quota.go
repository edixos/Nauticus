package controllers

import (
	"context"
	"fmt"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileResourceQuota(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Status.NamespaceName,
		},
		Spec: space.Spec.ResourceQuota,
	}
	if err := controllerutil.SetControllerReference(space, resourceQuota, s.Scheme); err != nil {
		return fmt.Errorf("unable to fill the ownerreference for the namespace")
	}
	existingResourceQuota := &corev1.ResourceQuota{}
	err := s.Client.Get(ctx, client.ObjectKey{Name: space.Name, Namespace: space.Status.NamespaceName}, existingResourceQuota)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Creating the resource quota", "name", space.Name, "namespace", space.Status.NamespaceName)
			err = s.Client.Create(ctx, resourceQuota)
			if err != nil {
				log.Error(err, "Failed to create resource quota", "name", space.Name)
				return err
			}
			log.Info("The resource-quota created successfully.", "name", space.Name)
		} else {
			// update existing resource quota
			log.Info("Resource quota already existing, updating it.")
			err = s.Client.Update(ctx, resourceQuota)
		}
	}
	return nil

}
