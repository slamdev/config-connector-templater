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
	gke "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	api "github.com/slamdev/config-connector-templater/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	"time"
)

var k8sClient client.Client

func TestTemplateReconciler(t *testing.T) {
	ctx := context.Background()

	stringPtr := func(s string) *string {
		return &s
	}

	const (
		TemplateResName = "test-template"
		Namespace       = "default"

		TemplatedResID = "{{ .metadata.namespace }}.test-ref"
		RenderedResID  = Namespace + ".test-ref"

		TemplatedRegion = "{{ .metadata.namespace }}.us"
		RenderedRegion  = Namespace + ".us"

		TemplatedKeyName = "{{ .metadata.namespace }}.kms"
		RenderedKeyName  = Namespace + ".kms"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	res := &api.PubSubTopicTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "config-connector-templater.slamdev.net/v1alpha1",
			Kind:       "PubSubTopicTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TemplateResName,
			Namespace: Namespace,
		},
		Spec: pubsub.PubSubTopicSpec{
			ResourceID: stringPtr(TemplatedResID),
			MessageStoragePolicy: &pubsub.TopicMessageStoragePolicy{
				AllowedPersistenceRegions: []string{TemplatedRegion},
			},
			KmsKeyRef: &gke.ResourceRef{
				Name:      TemplatedKeyName,
				Namespace: Namespace,
			},
		},
	}

	assert.NoError(t, k8sClient.Create(ctx, res))

	lookupKey := types.NamespacedName{Name: TemplateResName, Namespace: Namespace}
	createdRes := &api.PubSubTopicTemplate{}

	assert.Eventually(t, func() bool {
		err := k8sClient.Get(ctx, lookupKey, createdRes)
		if err != nil {
			return false
		}
		return true
	}, timeout, interval)

	assert.Equal(t, res.Spec.ResourceID, createdRes.Spec.ResourceID)

	renderedRes := &pubsub.PubSubTopic{}

	assert.Eventually(t, func() bool {
		if err := k8sClient.Get(ctx, lookupKey, renderedRes); err != nil {
			return false
		}
		return true
	}, timeout, interval)

	assert.Equal(t, RenderedResID, *renderedRes.Spec.ResourceID)
	assert.Equal(t, RenderedRegion, renderedRes.Spec.MessageStoragePolicy.AllowedPersistenceRegions[0])
	assert.Equal(t, RenderedKeyName, renderedRes.Spec.KmsKeyRef.Name)
	assert.Equal(t, Namespace, renderedRes.Spec.KmsKeyRef.Namespace)

	assert.NoError(t, k8sClient.Get(ctx, lookupKey, res))
	res.Spec.MessageStoragePolicy = &pubsub.TopicMessageStoragePolicy{AllowedPersistenceRegions: []string{"test"}}
	assert.NoError(t, k8sClient.Update(ctx, res))

	assert.Eventually(t, func() bool {
		if err := k8sClient.Get(ctx, lookupKey, renderedRes); err != nil {
			return false
		}
		return renderedRes.Spec.MessageStoragePolicy != nil
	}, timeout, interval)

	assert.Equal(t, res.Spec.MessageStoragePolicy.AllowedPersistenceRegions, renderedRes.Spec.MessageStoragePolicy.AllowedPersistenceRegions)
}

func TestMain(m *testing.M) {
	// setUp
	if os.Getenv("KUBEBUILDER_ASSETS") == "" && os.Getenv("ENVTEST_ASSETS_DIR") == "" {
		// allows to run tests from GoLand
		_ = os.Setenv("KUBEBUILDER_ASSETS", "../testbin/bin")
	}
	logf.SetLogger(zap.New(zap.WriteTo(os.Stdout), zap.UseDevMode(true)))
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
			filepath.Join("..", "testbin", "crd"),
		},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	if err != nil {
		panic(err)
	}

	if err := pubsub.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	if err := api.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		panic(err)
	}

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	if err != nil {
		panic(err)
	}

	if err := CreateControllers(k8sManager); err != nil {
		panic(err)
	}

	go func() {
		if err := k8sManager.Start(ctrl.SetupSignalHandler()); err != nil {
			panic(err)
		}
	}()

	code := m.Run()

	// tearDown
	if err := testEnv.Stop(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
