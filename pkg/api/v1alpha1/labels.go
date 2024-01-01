// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"fmt"

	"github.com/edixos/nauticus/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTypeLabel(t metav1.Object) (label string, err error) {
	switch v := t.(type) {
	case *v1alpha1.Space:
		return "nauticus.io/space", nil
	case *corev1.Namespace:
		return "nauticus.io/namespace", nil
	case *networkingv1.NetworkPolicy:
		return "nauticus.io/network-policy", nil
	case *corev1.ResourceQuota:
		return "nauticus.io/resource-quota", nil
	case *rbacv1.RoleBinding:
		return "nauticus.io/role-binding", nil
	case *corev1.LimitRange:
		return "nauticus.io/limit-range", nil
	case *corev1.ServiceAccount:
		return "nauticus.io/service-account", nil
	default:
		err = fmt.Errorf("type %T is not mapped as Nauticus label recognized", v)
	}

	return
}
