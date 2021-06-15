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

package pkg

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CliCli interface {
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
	Status() client.StatusWriter
	GetScheme() *runtime.Scheme
}

func CreateTargetResource(ctx context.Context, cli CliCli, src client.Object, typedContainer client.Object) error {
	if err := createTemplatedResource(cli, src, typedContainer); err != nil {
		return fmt.Errorf("failed to create templated resource; %w", err)
	}
	return cli.Create(ctx, typedContainer)
}

func UpdateTargetResource(ctx context.Context, cli CliCli, src client.Object, target client.Object, typedContainer client.Object) error {
	// Build the PubSubTopic spec from PubSubTopicTemplate
	if err := createTemplatedResource(cli, src, typedContainer); err != nil {
		return fmt.Errorf("failed to create templated resource; %w", err)
	}
	resSpec := getSpec(typedContainer)

	// Update PubSubTopic if needed
	if !reflect.DeepEqual(resSpec, getSpec(target)) {
		setSpec(target, resSpec)
		err := cli.Update(ctx, target)
		if err != nil {
			return fmt.Errorf("failed to update dest resource; %w", err)
		}
	}

	if err := updateStatusRef(ctx, cli, target, src); err != nil {
		return fmt.Errorf("failed to update dest resource status; %w", err)
	}

	return nil
}

func updateStatusRef(ctx context.Context, cli CliCli, src client.Object, target client.Object) error {
	// Build the PubSubTopicTemplate status ref with from PubSubTopic
	ref := corev1.ObjectReference{
		Kind:            src.GetObjectKind().GroupVersionKind().Kind,
		Namespace:       src.GetNamespace(),
		Name:            src.GetName(),
		UID:             src.GetUID(),
		APIVersion:      src.GetObjectKind().GroupVersionKind().Version,
		ResourceVersion: src.GetResourceVersion(),
	}

	// Update status.Ref if needed
	if !reflect.DeepEqual(ref, getStatusRef(target)) {
		setStatusRef(target, ref)
		return cli.Status().Update(ctx, target)
	}

	return nil
}

func createTemplatedResource(cli CliCli, src client.Object, target client.Object) error {
	spec, err := Render(getSpec(src), src)
	if err != nil {
		return fmt.Errorf("failed to render template; %w", err)
	}
	target.SetName(src.GetName())
	target.SetNamespace(src.GetNamespace())
	target.SetLabels(src.GetLabels())
	target.SetAnnotations(src.GetAnnotations())

	setSpec(target, spec)

	if err := ctrl.SetControllerReference(src, target, cli.GetScheme()); err != nil {
		return fmt.Errorf("failed to set ctrl ref; %w", err)
	}
	return nil
}

func setSpec(target interface{}, spec interface{}) {
	v := reflect.ValueOf(target)
	f := v.FieldByName("Spec")
	f.Set(reflect.ValueOf(spec))
}

func getSpec(target interface{}) interface{} {
	v := reflect.ValueOf(target)
	f := v.FieldByName("Spec")
	return f.Elem().Interface()
}

func setStatusRef(target interface{}, ref corev1.ObjectReference) {
	v := reflect.ValueOf(target)
	statusValue := v.FieldByName("Status")
	refValue := statusValue.FieldByName("Ref")
	refValue.Set(reflect.ValueOf(ref))
}

func getStatusRef(target interface{}) corev1.ObjectReference {
	v := reflect.ValueOf(target)
	statusValue := v.FieldByName("Status")
	refValue := statusValue.FieldByName("Ref")
	return refValue.Elem().Interface().(corev1.ObjectReference)
}
