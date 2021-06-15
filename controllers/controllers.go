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
	"fmt"
	pubsub "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/pubsub/v1beta1"
	api "github.com/slamdev/config-connector-templater/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubtopictemplates/finalizers,verbs=update
//+kubebuilder:rbac:groups=pubsub.cnrm.cloud.google.com,resources=pubsubtopics,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubsubscriptiontemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubsubscriptiontemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config-connector-templater.slamdev.net,resources=pubsubsubscriptiontemplates/finalizers,verbs=update
//+kubebuilder:rbac:groups=pubsub.cnrm.cloud.google.com,resources=pubsubsubscriptions,verbs=get;list;watch;create;update;patch;delete

var controlledTypes = []controlledType{
	{
		id:           "pubsubtopictemplate",
		templateType: &api.PubSubTopicTemplate{},
		renderType:   &pubsub.PubSubTopic{},
	},
	{
		id:           "pubsubsubscriptiontemplate",
		templateType: &api.PubSubSubscriptionTemplate{},
		renderType:   &pubsub.PubSubSubscription{},
	},
}

type controlledType struct {
	id           string
	templateType client.Object
	renderType   client.Object
}

func CreateControllers(mgr ctrl.Manager) error {
	for _, t := range controlledTypes {
		c := &TemplateReconciler{
			Client:       mgr.GetClient(),
			Scheme:       mgr.GetScheme(),
			LoggerName:   t.id,
			TemplateType: t.templateType,
			RenderType:   t.renderType,
		}
		if err := c.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create %s controller; %w", c.LoggerName, err)
		}
	}
	return nil
}
