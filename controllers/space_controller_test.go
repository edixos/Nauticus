// /*
// Copyright 2023.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */
package controllers

import (
	"context"
	"time"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.
const (
	SpaceName                  = "test-space"
	SpaceNameWithResourceQuota = "test-space-resource-quota"
	timeout                    = time.Second * 10
	interval                   = time.Millisecond * 250
)

var _ = Describe("Space controller", func() {
	Context("When creating a basic space", func() {
		var createdSpace nauticusiov1alpha1.Space

		It("Create a basic resource", func() {
			By("Creating a Space", func() {
				ctx := context.Background()
				space := &nauticusiov1alpha1.Space{
					TypeMeta: metav1.TypeMeta{
						APIVersion: nauticusiov1alpha1.GroupVersion.Version,
						Kind:       nauticusiov1alpha1.SpaceKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: SpaceName,
					},
				}
				Expect(k8sClient.Create(ctx, space)).Should(Succeed())
				spaceLookupKey := types.NamespacedName{Name: space.Name}
				// We'll need to retry getting this newly created Space, given that creation may not immediately happen.
				Eventually(func() error {
					return k8sClient.Get(ctx, spaceLookupKey, &createdSpace)
				}, timeout, interval).Should(Succeed())
			})
		})
		It("Should create a Namespace", func() {
			ctx := context.Background()
			By("Creating an instance with a generated name", func() {
				namespaceLookupKey := types.NamespacedName{Name: createdSpace.Status.NamespaceName}
				createdNamespace := &corev1.Namespace{}
				// We'll need to retry getting this newly created MsTeams, given that creation may not immediately happen.
				Eventually(func() error {
					return k8sClient.Get(ctx, namespaceLookupKey, createdNamespace)
				}, timeout, interval).Should(Succeed())
				Expect(createdNamespace.OwnerReferences[0].UID).To(Equal(createdSpace.UID))
				Expect(createdSpace.Status.NamespaceName).To(Equal(createdNamespace.Name))
			})
		})
	})

	Context("When creating a space with resource quota", func() {
		var createdSpaceWithQuota nauticusiov1alpha1.Space
		var createdResourceQuota corev1.ResourceQuota
		ctx := context.Background()
		It("Creates a Space with resource quota spec", func() {
			By("Creating a Space with resource quotas", func() {
				spaceWithResourceQuota := &nauticusiov1alpha1.Space{
					TypeMeta: metav1.TypeMeta{
						APIVersion: nauticusiov1alpha1.GroupVersion.Version,
						Kind:       nauticusiov1alpha1.SpaceKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: SpaceNameWithResourceQuota,
					},
					Spec: nauticusiov1alpha1.SpaceSpec{
						ResourceQuota: corev1.ResourceQuotaSpec{
							Hard: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("8"),
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, spaceWithResourceQuota)).Should(Succeed())
				spaceLookupKey := types.NamespacedName{Name: spaceWithResourceQuota.Name}
				// We'll need to retry getting this newly created Space, given that creation may not immediately happen.
				Eventually(func() error {
					return k8sClient.Get(ctx, spaceLookupKey, &createdSpaceWithQuota)
				}, timeout, interval).Should(Succeed())
			})
		})
		It("Should create a resource quota", func() {
			resourceQuotaLookupKey := types.NamespacedName{
				Namespace: createdSpaceWithQuota.Status.NamespaceName,
				Name:      createdSpaceWithQuota.Name,
			}
			By("Creating a resource quota within the generated namespace", func() {
				// We'll need to retry getting this newly created resource-quota, given that creation may not immediately happen.
				Eventually(func() error {
					return k8sClient.Get(ctx, resourceQuotaLookupKey, &createdResourceQuota)
				}, timeout, interval).Should(Succeed())
			})
		})
	})
})
