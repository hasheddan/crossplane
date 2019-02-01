/*
Copyright 2018 The Crossplane Authors.

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

package v1alpha1

import (
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProviderSpec defines the desired state of Provider
type ProviderSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file

	// Azure service principal credentials json secret key reference
	Secret corev1.SecretKeySelector `json:"credentialsSecretRef"`
}

// ProviderStatus is the status for this provider
type ProviderStatus struct {
	corev1alpha1.ConditionedStatus
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Provider is the Schema for the instances API
// +k8s:openapi-gen=true
// +groupName=azure
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec   `json:"spec,omitempty"`
	Status ProviderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

// IsValid returns true if provider is valid (in ready state)
func (p *Provider) IsValid() bool {
	return p.Status.IsReady()
}

// ResourceGroupSpec defines the desired state of ResourceGroup
// https://godoc.org/github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources#Group
type ResourceGroupSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file

	// ID - The ID of the resource group.
	ID string `json:"id,omitempty"`

	// Name - The name of the resource group.
	// +kubebuilder:validation:MaxLength=90
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// Type - The type of the resource group.
	Type string `json:"type,omitempty"`

	// Properties GroupProperties `json:"properties,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	// +kubebuilder:validation:MinLength=1
	Location string `json:"location,omitempty"`

	// ManagedBy - The ID of the resource that manages this resource group.
	ManagedBy string `json:"managedBy,omitempty"`

	// Tags - The tags attached to the resource group.
	Tags map[string]string `json:"tags"`
}

// ResourceGroupStatus is the status for this ResourceGroup
type ResourceGroupStatus struct {
	corev1alpha1.ConditionedStatus
	// Rationale here is that we do not want the reconciler to know when to create a new Resource Group, so if the name is not set we know it needs to be created
	Name string `json:"name,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceGroup is the Schema for the ResourceGroup API
// +k8s:openapi-gen=true
// +groupName=azure
type ResourceGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceGroupSpec   `json:"spec,omitempty"`
	Status ResourceGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProviderList contains a list of Provider
type ResourceGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceGroup `json:"items"`
}
