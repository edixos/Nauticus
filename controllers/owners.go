package controllers

import (
    "context"
    "fmt"

    nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
    "github.com/go-logr/logr"
    rbacv1 "k8s.io/api/rbac/v1"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileOwners(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
    rolebindingName := space.Name + "-owner"
    ownersRoleBinding := &rbacv1.RoleBinding{
        ObjectMeta: metav1.ObjectMeta{
            Name:      rolebindingName,
            Namespace: space.Status.NamespaceName,
        },
        RoleRef: rbacv1.RoleRef{
            Kind:     "ClusterRole",
            APIGroup: rbacv1.GroupName,
            Name:     "admin",
        },
        Subjects: space.Spec.Owners,
    }
    if err := controllerutil.SetControllerReference(space, ownersRoleBinding, s.Scheme); err != nil {
        return fmt.Errorf("unable to fill the ownerreference for the owners rolebindings")
    }

    existingRoleBinding := &rbacv1.RoleBinding{}
    roleBindingLookupKey := types.NamespacedName{Name: rolebindingName, Namespace: space.Status.NamespaceName}
    err := s.Client.Get(ctx, roleBindingLookupKey, existingRoleBinding)
    if err != nil {
        if apierrors.IsNotFound(err) {
            log.Info("Creating the role binding", "name", space.Name, "namespace", space.Status.NamespaceName)
            err = s.Client.Create(ctx, ownersRoleBinding)
            if err != nil {
                log.Error(err, "Failed to create role binding", "name", space.Name)
                return err
            }
            log.Info("The role binding created successfully.", "name", space.Name)
        } else {
            // update existing resource quota
            log.Info("Unable to fetch if there is an existing role binding.")
            return err
        }
    } else {
        log.Info("Updating existing RoleBinding", "RoleBinding", roleBindingLookupKey)
        err = s.Client.Update(ctx, ownersRoleBinding)
        if err != nil {
            log.Error(err, "Cannot Update the existing Rolebinding", "RoleBinding", roleBindingLookupKey)
        }
    }
    return nil

}
