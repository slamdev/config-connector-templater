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
	"fmt"
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	"github.com/Masterminds/sprig/v3"
	api "github.com/slamdev/config-connector-templater/api/v1alpha1"
	"html/template"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

// PubSubTopicTemplateReconciler reconciles a PubSubTopicTemplate object
type PubSubTopicTemplateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates/finalizers,verbs=update
//+kubebuilder:rbac:groups=pubsub.cnrm.cloud.google.com,resources=pubsubtopics,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PubSubTopicTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("pubsubtopictemplate", req.NamespacedName)

	// Fetch the Memcached instance
	templateRes := &api.PubSubTopicTemplate{}
	err := r.Get(ctx, req.NamespacedName, templateRes)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("PubSubTopicTemplate resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get PubSubTopicTemplate")
		return ctrl.Result{}, err
	}

	// Check if the PubSubTopic already exists, if not create a new one
	found := &pubsub.PubSubTopic{}
	err = r.Get(ctx, types.NamespacedName{Name: templateRes.Name, Namespace: templateRes.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new PubSubTopic resource
		res, err := r.createTemplatedResource(templateRes)
		if err != nil {
			logger.Error(err, "Failed to create templated resource")
			return ctrl.Result{}, err
		}
		logger.Info("Creating a new PubSubTopic", "PubSubTopic.Namespace", res.Namespace, "PubSubTopic.Name", res.Name)
		err = r.Create(ctx, res)
		if err != nil {
			logger.Error(err, "Failed to create new PubSubTopic", "PubSubTopic.Namespace", res.Namespace, "PubSubTopic.Name", res.Name)
			return ctrl.Result{}, err
		}
		// PubSubTopic created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get PubSubTopic")
		return ctrl.Result{}, err
	}

	// Build the PubSubTopic spec from PubSubTopicTemplate
	res, err := r.createTemplatedResource(templateRes)
	if err != nil {
		logger.Error(err, "Failed to create templated resource")
		return ctrl.Result{}, err
	}
	resSpec := res.Spec

	// Update PubSubTopic if needed
	if !reflect.DeepEqual(resSpec, found.Spec) {
		found.Spec = resSpec
		err := r.Status().Update(ctx, found)
		if err != nil {
			logger.Error(err, "Failed to update PubSubTopic spec")
			return ctrl.Result{}, err
		}
	}

	// Build the PubSubTopicTemplate status ref with from PubSubTopic
	ref := corev1.ObjectReference{
		Kind:            found.Kind,
		Namespace:       found.Namespace,
		Name:            found.Name,
		UID:             found.UID,
		APIVersion:      found.APIVersion,
		ResourceVersion: found.ResourceVersion,
	}

	// Update status.Ref if needed
	if !reflect.DeepEqual(ref, templateRes.Status.Ref) {
		templateRes.Status.Ref = ref
		err := r.Status().Update(ctx, templateRes)
		if err != nil {
			logger.Error(err, "Failed to update PubSubTopicTemplate status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *PubSubTopicTemplateReconciler) createTemplatedResource(template *api.PubSubTopicTemplate) (*pubsub.PubSubTopic, error) {
	spec, err := r.render(template.Spec, template)
	if err != nil {
		return nil, fmt.Errorf("failed to render template; %w", err)
	}
	res := &pubsub.PubSubTopic{
		ObjectMeta: metav1.ObjectMeta{
			Name:        template.Name,
			Namespace:   template.Namespace,
			Labels:      template.Labels,
			Annotations: template.Annotations,
		},
		Spec: spec.(pubsub.PubSubTopicSpec),
	}
	if err := ctrl.SetControllerReference(template, res, r.Scheme); err != nil {
		return nil, fmt.Errorf("failed to set ctrl ref; %w", err)
	}
	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PubSubTopicTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.PubSubTopicTemplate{}).
		Owns(&pubsub.PubSubTopic{}).
		Complete(r)
}

func (r *PubSubTopicTemplateReconciler) render(spec interface{}, data interface{}) (interface{}, error) {
	if err := r.renderValue(reflect.ValueOf(spec).Elem(), data); err != nil {
		return nil, fmt.Errorf("failed to render spec; %w", err)
	}
	return spec, nil
}

func (r *PubSubTopicTemplateReconciler) renderValue(v reflect.Value, data interface{}) error {
	s := v
	tpl := template.New("_").Funcs(sprig.FuncMap())
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Kind() == reflect.String && f.CanSet() {
			parsed, err := tpl.Parse(f.String())
			if err != nil {
				return fmt.Errorf("failed to parse template; %w", err)
			}
			v := new(strings.Builder)
			if err := parsed.Execute(v, data); err != nil {
				return fmt.Errorf("failed to execute template; %w", err)
			}
			f.SetString(v.String())
		} else if f.Kind() == reflect.Struct {
			x1 := reflect.ValueOf(f.Interface())
			if err := r.renderValue(x1, data); err != nil {
				return fmt.Errorf("failed to render value; %w", err)
			}
		}
	}
	return nil
}
