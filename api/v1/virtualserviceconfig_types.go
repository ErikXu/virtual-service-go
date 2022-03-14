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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualServiceConfigSpec defines the desired state of VirtualServiceConfig
type VirtualServiceConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	VirtualServiceName string `json:"virtualServiceName"`

	// +kubebuilder:validation:Required
	Host string `json:"host"`

	// +kubebuilder:validation:Required
	Http []HttpRoute `json:"http"`
}

type HttpRoute struct {

	// +nullable
	Name string `json:"name,omitempty"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=0
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:ExclusiveMaximum=false
	// +kubebuilder:validation:Maximum=9999
	Order int32 `json:"order,omitempty"`

	// +kubebuilder:validation:Required
	Match HttpMatchRequest `json:"match"`

	// +kubebuilder:validation:Required
	Route HttpRouteDestination `json:"route"`
}

type HttpMatchRequest struct {

	// +nullable
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:Required
	Uri StringMatch `json:"uri"`

	// +kubebuilder:validation:Optional
	Headers map[string]StringMatch `json:"headers"`
}

type HttpRouteDestination struct {

	// +kubebuilder:validation:Required
	Host string `json:"host"`

	// +nullable
	Subset string `json:"subset,omitempty"`
}

type StringMatch struct {
	Exact string `json:"exact,omitempty"`

	Prefix string `json:"prefix,omitempty"`

	Regex string `json:"regex,omitempty"`
}

// VirtualServiceConfigStatus defines the observed state of VirtualServiceConfig
type VirtualServiceConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=vsc

// VirtualServiceConfig is the Schema for the virtualserviceconfigs API
type VirtualServiceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualServiceConfigSpec   `json:"spec,omitempty"`
	Status VirtualServiceConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VirtualServiceConfigList contains a list of VirtualServiceConfig
type VirtualServiceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualServiceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualServiceConfig{}, &VirtualServiceConfigList{})
}
