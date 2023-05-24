package controllers

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileOwners(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	rolebindingName := space.Name + "-owner"

	roleRef := rbacv1.RoleRef{
		Kind:     "ClusterRole",
		APIGroup: rbacv1.GroupName,
		Name:     "admin",
	}
	ownersRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, roleRef, space.Spec.Owners)

	err = s.syncRoleBinding(ctx, ownersRoleBinding, space, ownersRoleBinding.RoleRef, ownersRoleBinding.Subjects)

	return err
}

func (s *SpaceReconciler) reconcileAdditionalRoleBindings(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	for _, ad := range space.Spec.AdditionalRoleBindings {
		rolebindingName := space.Name + "-" + ad.RoleRef.Name
		additionalRoleBinding := newRoleBinding(rolebindingName, space.Status.NamespaceName, ad.RoleRef, ad.Subjects)

		err = s.syncRoleBinding(ctx, additionalRoleBinding, space, ad.RoleRef, ad.Subjects)
	}

	return err
}

func (s *SpaceReconciler) syncRoleBinding(ctx context.Context, roleBinding *rbacv1.RoleBinding, space *nauticusiov1alpha1.Space, desiredRoleRef rbacv1.RoleRef, desiredSubjects []rbacv1.Subject) (err error) {
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

	res, err = controllerutil.CreateOrUpdate(ctx, s.Client, roleBinding, func() error {
		roleBinding.SetLabels(map[string]string{
			spaceLabel:       space.Name,
			roleBindingLabel: roleBinding.Name,
		})
		roleBinding.RoleRef = desiredRoleRef
		roleBinding.Subjects = desiredSubjects

		return nil
	})

	s.Log.Info("Rolebinding sync result: "+string(res), "name", roleBinding.Name, "namespace", space.Status.NamespaceName)
	s.emitEvent(space, space.Name, res, "Ensuring RoleBinding creation/Update", err)

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
