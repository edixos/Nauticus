// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package space

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) reconcileOwners(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	rolebindingName := space.Name + "-owner"

	roleRef := rbacv1.RoleRef{
		Kind:     "ClusterRole",
		APIGroup: rbacv1.GroupName,
		Name:     "admin",
	}
	ownersRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, roleRef, space.Spec.Owners)

	err = r.syncRoleBinding(ctx, ownersRoleBinding, space, ownersRoleBinding.RoleRef, ownersRoleBinding.Subjects)

	return err
}

func (r *Reconciler) reconcileAdditionalRoleBindings(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for _, ad := range space.Spec.AdditionalRoleBindings {
		rolebindingName := space.Name + "-" + ad.RoleRef.Name
		additionalRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, ad.RoleRef, ad.Subjects)

		err = r.syncRoleBinding(ctx, additionalRoleBinding, space, ad.RoleRef, ad.Subjects)
	}

	return err
}

func (r *Reconciler) syncRoleBinding(ctx context.Context, roleBinding *rbacv1.RoleBinding, space *nauticusiov1alpha1.Space, desiredRoleRef rbacv1.RoleRef, desiredSubjects []rbacv1.Subject) (err error) {
	var (
		res                          controllerutil.OperationResult
		spaceLabel, roleBindingLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if roleBindingLabel, err = v1alpha1.GetTypeLabel(roleBinding); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, roleBinding, func() error {
		roleBinding.SetLabels(map[string]string{
			spaceLabel:       space.Name,
			roleBindingLabel: roleBinding.Name,
		})
		roleBinding.RoleRef = desiredRoleRef
		roleBinding.Subjects = desiredSubjects

		return nil
	})

	r.Log.Info("Rolebinding sync result: "+string(res), "name", roleBinding.Name, "namespace", space.Status.NamespaceName)
	r.EmitEvent(space, space.Name, res, "Ensuring RoleBinding creation/Update", err)

	return err
}

func newRoleBinding(name string, namespace string, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}

func (r *Reconciler) deleteOwners(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	rolebindingName := space.Name + "-owner"

	roleRef := rbacv1.RoleRef{}
	ownersRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, roleRef, space.Spec.Owners)

	err = r.DeleteObject(ctx, ownersRoleBinding)

	return err
}

func (r *Reconciler) deleteAdditionalRoleBindings(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for _, ad := range space.Spec.AdditionalRoleBindings {
		rolebindingName := space.Name + "-" + ad.RoleRef.Name
		additionalRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, ad.RoleRef, ad.Subjects)

		err = r.DeleteObject(ctx, additionalRoleBinding)
	}

	return err
}
