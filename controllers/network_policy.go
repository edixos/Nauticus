// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"
	"strconv"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SpaceReconciler) reconcileNetworkPolicies(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	if space.Spec.NetworkPolicies.EnableDefaultStrictMode {
		networkPolicyName := fmt.Sprintf("nauticus-%s", space.Name)
		networkPolicySpec := newNetworkPolicyDefaultSpec()
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicySpec)
		err = s.syncNetworkPolicy(ctx, networkPolicy, space, networkPolicySpec)

		if err != nil {
			s.Log.Error(err, "Cannot Synchronize Network policy")

			return err
		}
	}

	for i, networkPolicy := range space.Spec.NetworkPolicies.Items {
		npName := "nauticus-custom-" + strconv.Itoa(i)
		np := newNetworkPolicy(npName, space.Status.NamespaceName, networkPolicy)
		err = s.syncNetworkPolicy(ctx, np, space, networkPolicy)

		if err != nil {
			s.Log.Error(err, "Cannot Synchronize Network policy")

			return err
		}
	}

	return nil
}

func (s *SpaceReconciler) syncNetworkPolicy(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, space *nauticusiov1alpha1.Space, spec networkingv1.NetworkPolicySpec) (err error) {
	var (
		res                            controllerutil.OperationResult
		spaceLabel, networkPolicyLabel string
	)

	if spaceLabel, err = v1alpha1.GetTypeLabel(space); err != nil {
		return
	}

	if networkPolicyLabel, err = v1alpha1.GetTypeLabel(networkPolicy); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, s.Client, networkPolicy, func() (err error) {
		networkPolicy.SetLabels(map[string]string{
			spaceLabel:         space.Name,
			networkPolicyLabel: networkPolicy.Name,
		})
		if networkPolicy.Name != fmt.Sprintf("nauticus-%s", space.Name) {
			networkPolicy.Spec = spec
		}

		return nil
	})
	s.Log.Info("Network Policy sync result: "+string(res), "name", networkPolicy.Name, "namespace", space.Status.NamespaceName)
	s.emitEvent(space, space.Name, res, "Ensuring NetworkPolicy creation/Update", err)

	return nil
}

func newNetworkPolicy(name string, namespace string, networkPolicySpec networkingv1.NetworkPolicySpec) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkPolicySpec,
	}
}

func newNetworkPolicyDefaultSpec() networkingv1.NetworkPolicySpec {
	return networkingv1.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{MatchLabels: map[string]string{}},
		Ingress: []networkingv1.NetworkPolicyIngressRule{
			{
				From: []networkingv1.NetworkPolicyPeer{
					{
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"nauticus.io/role": "system",
							},
						},
					},
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{},
						},
					},
				},
			},
		},
	}
}

func (s *SpaceReconciler) deleteNetworkPolicies(ctx context.Context, space *nauticusiov1alpha1.Space) (err error) {
	if space.Spec.NetworkPolicies.EnableDefaultStrictMode {
		networkPolicyName := fmt.Sprintf("nauticus-%s", space.Name)
		networkPolicySpec := newNetworkPolicyDefaultSpec()
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicySpec)

		if err = s.deleteObject(ctx, networkPolicy); err != nil {
			return err
		}
	}

	for i, networkPolicy := range space.Spec.NetworkPolicies.Items {
		npName := "nauticus-custom-" + strconv.Itoa(i)
		np := newNetworkPolicy(npName, space.Status.NamespaceName, networkPolicy)

		if err = s.deleteObject(ctx, np); err != nil {
			return err
		}
	}

	return nil
}
