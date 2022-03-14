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
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	extensionv1 "virtual-service-go/api/v1"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"encoding/json"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	log.Info(fmt.Sprintf("Virtual service config: [%s] is reconciled", req.NamespacedName))

	var config extensionv1.VirtualServiceConfig
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		log.Info(fmt.Sprintf("Virtual service config: [%s] is deleted", req.NamespacedName))

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	finalizerName := "virtualserviceconfigs.extension.networking.istio.io/finalizer"

	if config.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&config, finalizerName) {
			log.Info(fmt.Sprintf("Add finalizer for virtual service config: [%s]", req.NamespacedName))

			controllerutil.AddFinalizer(&config, finalizerName)
			if err := r.Update(ctx, &config); err != nil {
				log.Error(err, "Add finalizer failed")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(&config, finalizerName) {
			log.Info(fmt.Sprintf("Handle virtual service before config: [%s] is being deleted", req.NamespacedName))
			var list extensionv1.VirtualServiceConfigList
			if err := r.List(ctx, &list, client.InNamespace(req.Namespace)); err != nil {
				log.Error(err, "Unable to list configs")
				return ctrl.Result{}, err
			}

			configs := filterConfigs(list, config.Spec.VirtualServiceName, config.Name)
			if len(configs) > 0 {
				generateVirtualService(ctx, log, config.Spec.VirtualServiceName, req.Namespace, config.Spec.Host, configs)
			} else {
				deleteVirtualService(ctx, log, req.Namespace, config.Spec.VirtualServiceName)
			}

			log.Info(fmt.Sprintf("Remove finalizer for virtual service config: [%s]", req.NamespacedName))
			controllerutil.RemoveFinalizer(&config, finalizerName)
			if err := r.Update(ctx, &config); err != nil {
				log.Error(err, "Remove finalizer failed")
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Your reconcile logic

	var list extensionv1.VirtualServiceConfigList
	if err := r.List(ctx, &list, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "Unable to list configs")
		return ctrl.Result{}, err
	}

	configs := filterConfigs(list, config.Spec.VirtualServiceName, "")
	generateVirtualService(ctx, log, config.Spec.VirtualServiceName, req.Namespace, config.Spec.Host, configs)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualServiceConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.VirtualServiceConfig{}).
		Complete(r)
}

func filterConfigs(list extensionv1.VirtualServiceConfigList, virtualName string, exceptName string) []extensionv1.VirtualServiceConfig {
	configs := make([]extensionv1.VirtualServiceConfig, 0)

	for _, item := range list.Items {
		if len(exceptName) > 0 {
			if item.Spec.VirtualServiceName == virtualName && item.Name != exceptName {
				configs = append(configs, item)
			}
		} else {
			if item.Spec.VirtualServiceName == virtualName {
				configs = append(configs, item)
			}
		}
	}

	return configs
}

func generateVirtualService(ctx context.Context, log logr.Logger, virtualServiceName string, namespace string, host string, configs []extensionv1.VirtualServiceConfig) {
	log.Info("Generate virtual service by configs")
	virtualService := configsToVirtualService(virtualServiceName, namespace, host, configs)

	d, err := yaml.Marshal(&virtualService)
	if err != nil {
		log.Error(err, "Unable to convert to yaml")
	}

	log.Info(fmt.Sprintf("Generated virtual service yaml:\n%s", string(d)))

	applyVirtualService(ctx, log, virtualService)
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

	specHttps := make([]extensionv1.HttpRoute, 0)

	for _, config := range configs {
		specHttps = append(specHttps, config.Spec.Http...)
	}

	// Sort spechttps
	sort.SliceStable(specHttps, func(a, b int) bool {

		// Sort by order desc first
		if specHttps[a].Order > specHttps[b].Order {
			return true
		}

		// Then by uri len desc
		uriLenA := len(stringMatchToMap(specHttps[a].Match.Uri))
		uriLenB := len(stringMatchToMap(specHttps[b].Match.Uri))
		if specHttps[a].Order == specHttps[b].Order && uriLenA > uriLenB {
			return true
		}

		// Then by header count desc
		headerCountA := len(specHttps[a].Match.Headers)
		headerCountB := len(specHttps[b].Match.Headers)
		if specHttps[a].Order == specHttps[b].Order && uriLenA == uriLenB && headerCountA > headerCountB {
			return true
		}

		return false
	})

	for _, specHttp := range specHttps {
		http := Http{
			Match: []Match{},
			Route: []Route{},
		}

		if len(specHttp.Name) > 0 {
			http.Name = specHttp.Name
		}

		match := Match{
			Uri: stringMatchToMap(specHttp.Match.Uri),
		}

		if len(specHttp.Match.Name) > 0 {
			match.Name = specHttp.Match.Name
		}

		if len(specHttp.Match.Headers) > 0 {
			match.Headers = map[string]map[string]string{}

			for key, value := range specHttp.Match.Headers {
				match.Headers[key] = stringMatchToMap(value)
			}
		}

		http.Match = append(http.Match, match)

		route := Route{
			Destination: Destination{
				Host: specHttp.Route.Host,
			},
		}

		if len(specHttp.Route.Subset) > 0 {
			route.Destination.Subset = specHttp.Route.Subset
		}

		http.Route = append(http.Route, route)

		virtualService.Spec.Http = append(virtualService.Spec.Http, http)
	}

	return virtualService
}

func stringMatchToMap(match extensionv1.StringMatch) map[string]string {
	if len(match.Exact) > 0 {
		return map[string]string{"exact": match.Exact}
	}

	if len(match.Prefix) > 0 {
		return map[string]string{"prefix": match.Prefix}
	}

	return map[string]string{"regex": match.Regex}
}

func applyVirtualService(ctx context.Context, log logr.Logger, virtualService VirtualService) {
	log.Info(fmt.Sprintf("Apply virtual service: [%s/%s]", virtualService.Metadata.Namespace, virtualService.Metadata.Name))

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Unable to get kube config")
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Error(err, "Unable to init client")
	}

	var gvr = schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}

	body, err := json.Marshal(virtualService)
	if err != nil {
		log.Error(err, "Marshal failed")
	}

	utd, err := client.Resource(gvr).Namespace("default").Patch(ctx, virtualService.Metadata.Name, types.ApplyPatchType, body, metav1.PatchOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		log.Error(err, "Apply failed")
	}

	result, err := utd.MarshalJSON()
	if err != nil {
		log.Error(err, "Marshal json failed")
	}

	log.Info(fmt.Sprintf("Virtual service: [%s/%s] is applied, apply response json:\n%s", virtualService.Metadata.Namespace, virtualService.Metadata.Name, string(result)))
}

func deleteVirtualService(ctx context.Context, log logr.Logger, namespace string, name string) {
	log.Info(fmt.Sprintf("Delete virtual service: [%s/%s]", namespace, name))
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Unable to get kube config")
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Error(err, "Unable to init client")
	}

	var gvr = schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}

	if err := client.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		log.Error(err, fmt.Sprintf("Delete virtual service: [%s/%s] failed", namespace, name))
	}

	log.Info(fmt.Sprintf("Virtual service: [%s/%s] is deleted", namespace, name))
}

type VirtualService struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type Spec struct {
	Hosts []string `json:"hosts" yaml:"hosts"`
	Http  []Http   `json:"http" yaml:"http"`
}

type Http struct {
	Name  string  `json:"name,omitempty" yaml:"name,omitempty"`
	Match []Match `json:"match" yaml:"match"`
	Route []Route `json:"route" yaml:"route"`
}

type Match struct {
	Name    string                       `json:"name,omitempty" yaml:"name,omitempty"`
	Uri     map[string]string            `json:"uri" yaml:"uri"`
	Headers map[string]map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type Route struct {
	Destination Destination `json:"destination" yaml:"destination"`
}

type Destination struct {
	Host   string `json:"host" yaml:"host"`
	Subset string `json:"subset,omitempty" yaml:"subset,omitempty"`
}
