///*
//Copyright 2023.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//*/
//
package controllers

import (
    "context"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    //+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var _ = Describe("Space controller", func() {
    const (
        SpaceName                  = "test-space"
        SpaceNameWithResourceQuota = "test-space-resource-quota"
        timeout                    = time.Second * 10
        interval                   = time.Millisecond * 250
    )
    Context("When creating a space", func() {
        It("Should adds the Status.NamespaceName", func() {
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
            err := k8sClient.Create(ctx, space)
            Expect(err).NotTo(HaveOccurred())

            createdSpace := &nauticusiov1alpha1.Space{}
            err = k8sClient.Get(ctx, client.ObjectKey{
                Name: SpaceName,
            }, createdSpace)
            Expect(createdSpace.Status.NamespaceName).ToNot(BeNil())
        })
    })
    Context("When creating a space with a resource quota", func() {
        It("Should create a resource quota within the generated namespace", func() {
            ctx := context.TODO()

            space := &nauticusiov1alpha1.Space{
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
            err := k8sClient.Create(ctx, space)
            Expect(err).NotTo(HaveOccurred())
            createdResourceQuota := &corev1.ResourceQuota{}
            createdSpace := &nauticusiov1alpha1.Space{}
            err = k8sClient.Get(ctx, client.ObjectKey{
                Name: space.Name,
            }, createdSpace)
            Expect(err).NotTo(HaveOccurred())

            // We'll need to retry getting this newly created space, given that creation may not immediately happen.

            Eventually(func() error {
                return k8sClient.Get(ctx, client.ObjectKey{
                    Name:      space.Name,
                    Namespace: space.Status.NamespaceName,
                }, createdResourceQuota)
            }, timeout, interval).Should(Succeed())

        })
    })
})
