package controllers

import (
	"context"
	"fmt"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileNamespace(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
	namespace, err := s.newNamespace(space)
	if err != nil {
		return err
	}
	// Check if the namespace exist and create it if it does not
	existingNamespace := &v1.Namespace{}
	err = s.Client.Get(ctx, client.ObjectKey{Name: space.Status.NamespaceName}, existingNamespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Namespace does not exist, creating it.
			log.Info("Creating the namespace", "namespace", namespace.Name)
			err = s.Client.Create(ctx, namespace)
			if err != nil {
				log.Error(err, "Failed to create namespace", "namespace", namespace.Name)
				return err
			}
			log.Info("The namespace created successfully.", "namespace", namespace.Name)
		} else {
			log.Error(err, "Failed to check if the namespace exists", "namespace", namespace.Name)
		}
	} else {
		namespace = existingNamespace
	}
	// Update the Space's status
	space.Status.NamespaceName = namespace.Name
	err = s.Client.Status().Update(ctx, space)
	if err != nil {
		log.Error(err, "Failed to update Space status", "space", space.Name)
		return err
	}
	return nil
}

func (s *SpaceReconciler) newNamespace(space *nauticusiov1alpha1.Space) (*v1.Namespace, error) {

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.namespaceName(space),
			Labels: space.Labels,
		},
	}
	if err := controllerutil.SetControllerReference(space, namespace, s.Scheme); err != nil {
		return nil, fmt.Errorf("unable to fill the ownerreference for the namespace")
	}
	return namespace, nil
}

func (s *SpaceReconciler) namespaceName(space *nauticusiov1alpha1.Space) string {
	return space.Name
}
