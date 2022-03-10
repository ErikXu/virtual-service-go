/*
Copyright 2022.

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
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	extensionv1 "virtual-service-go/api/v1"

	"gopkg.in/yaml.v2"

	linq "github.com/ahmetb/go-linq/v3"
)

// VirtualServiceConfigReconciler reconciles a VirtualServiceConfig object
type VirtualServiceConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=extension.networking.istio.io,resources=virtualserviceconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.networking.istio.io,resources=virtualserviceconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.networking.istio.io,resources=virtualserviceconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtualServiceConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *VirtualServiceConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var config extensionv1.VirtualServiceConfig
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		log.Info(fmt.Sprintf("Virtual service config: [%s] deleted", req.NamespacedName))

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info(fmt.Sprintf("Virtual service config: [%s] reconciled", req.NamespacedName))

	var configList extensionv1.VirtualServiceConfigList
	if err := r.List(ctx, &configList); err != nil {
		log.Error(err, "Unable to list configs。")
		return ctrl.Result{}, err
	}

	if len(configList.Items) <= 0 {
		err := errors.New("config list items are empty")
		log.Error(err, "Config list items are empty")
		return ctrl.Result{}, err
	}

	generateVirtualService(log, config.Spec.VirtualServiceName, config.Namespace, config.Spec.Host, configList.Items)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualServiceConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.VirtualServiceConfig{}).
		Complete(r)
}

func generateVirtualService(log logr.Logger, virtualServiceName string, namespace string, host string, configs []extensionv1.VirtualServiceConfig) {
	virtualService := configsToVirtualService(virtualServiceName, namespace, host, configs)

	d, err := yaml.Marshal(&virtualService)
	if err != nil {
		log.Error(err, "Unable to convert to yaml")
	}

	log.Info(fmt.Sprintf("Generated yaml:\n%s\n\n", string(d)))
}

func configsToVirtualService(virtualServiceName string, namespace string, host string, configs []extensionv1.VirtualServiceConfig) VirtualService {
	virtualService := VirtualService{
		ApiVersion: "networking.istio.io/v1alpha3",
		Kind:       "VirtualService",
		Metadata: Metadata{
			Namespace: namespace,
			Name:      virtualServiceName,
		},
		Spec: Spec{
			Hosts: []string{host},
			Http:  []Http{},
		},
	}

	var specHttps []extensionv1.HttpRoute

	linq.From(configs).SelectMany(
		func(config interface{}) linq.Query {
			return linq.From(config.(extensionv1.VirtualServiceConfig).Spec.Http)
		}).ToSlice(&specHttps)

	fmt.Printf("Http count: %d", len(specHttps))
	return virtualService
}

type VirtualService struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Spec struct {
	Hosts []string `json:"hosts"`
	Http  []Http   `json:"http"`
}

type Http struct {
	Name string `json:"name,omitempty"`
}
