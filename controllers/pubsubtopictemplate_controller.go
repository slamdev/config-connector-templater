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
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configconnectortemplaterv1alpha1 "github.com/slamdev/config-connector-templater/api/v1alpha1"
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PubSubTopicTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PubSubTopicTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PubSubTopicTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configconnectortemplaterv1alpha1.PubSubTopicTemplate{}).
		Owns(&v1beta1.PubSubTopic{}).
		Complete(r)
}
