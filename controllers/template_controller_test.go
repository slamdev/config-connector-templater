/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/slamdev/config-connector-templater/api/v1alpha1"
)

var _ = Describe("Template controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		TemplateResName = "test-template"
		Namespace       = "default"
		Ref             = "{{ .metadata.namespace }}.test-ref"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating CronJob Status", func() {
		It("Should increase CronJob Status.Active count when new Jobs are created", func() {
			By("By creating a new CronJob")
			ctx := context.Background()

			stringPtr := func(s string) *string {
				return &s
			}

			cronJob := &api.PubSubTopicTemplate{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "config-connector-templater.slamdev.net/v1alpha1",
					Kind:       "PubSubTopicTemplate",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      TemplateResName,
					Namespace: Namespace,
				},
				Spec: pubsub.PubSubTopicSpec{
					ResourceID: stringPtr(Ref),
				},
			}
			Expect(k8sClient.Create(ctx, cronJob)).Should(Succeed())

			cronjobLookupKey := types.NamespacedName{Name: TemplateResName, Namespace: Namespace}
			createdCronjob := &api.PubSubTopicTemplate{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdCronjob.Spec.ResourceID).Should(Equal(stringPtr(Ref)))

			//
			By("By checking that the CronJob has one active Job")
			Eventually(func() (string, error) {
				res := &pubsub.PubSubTopic{}
				err := k8sClient.Get(ctx, cronjobLookupKey, res)
				if err != nil {
					return "", err
				}
				return *res.Spec.ResourceID, nil
			}, timeout, interval).Should(Equal(Ref), "should list our active job %s in the active jobs list in status", Ref)
		})
	})
})
