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
	"github.com/slamdev/config-connector-templater/pkg"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TemplateReconciler reconciles a PubSubTopicTemplate object
type TemplateReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	LoggerName   string
	TemplateType client.Object
	RenderType   client.Object
}

func (r *TemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues(r.LoggerName, req.NamespacedName)

	res := r.initTemplateType()
	err := r.Get(ctx, req.NamespacedName, res)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get resource")
		return ctrl.Result{}, err
	}

	found := r.initRenderType()
	err = r.Get(ctx, types.NamespacedName{Name: res.GetName(), Namespace: res.GetNamespace()}, found)

	if err != nil && errors.IsNotFound(err) {
		if err := pkg.CreateTargetResource(ctx, r, res, r.initRenderType()); err != nil {
			logger.Error(err, "Failed to create resource")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get resource")
		return ctrl.Result{}, err
	}

	if err := pkg.UpdateTargetResource(ctx, r, res, found, r.initRenderType()); err != nil {
		logger.Error(err, "Failed to update resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TemplateReconciler) initTemplateType() client.Object {
	return reflect.New(reflect.ValueOf(r.TemplateType).Elem().Type()).Interface().(client.Object)
}

func (r *TemplateReconciler) initRenderType() client.Object {
	return reflect.New(reflect.ValueOf(r.RenderType).Elem().Type()).Interface().(client.Object)
}

// SetupWithManager sets up the controller with the Manager.
func (r *TemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.initTemplateType()).
		Owns(r.initRenderType()).
		Complete(r)
}

func (r *TemplateReconciler) GetScheme() *runtime.Scheme {
	return r.Scheme
}
