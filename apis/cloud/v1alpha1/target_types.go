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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ServerEndpoint struct {
	ClientCIDR    string `json:"clientCIDR,omitempty"`
	ServerAddress string `json:"serverAddress"`
}

type MasterAuth struct {
	ClusterCACertificate string `json:"clusterCACertificate,omitempty"`
}

// TargetSpec defines the desired state of Target
type TargetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Target. Edit target_types.go to remove/update
	ServerEndpoints       []ServerEndpoint `json:"serverEndpoints,omitempty"`
	MasterAuth            *MasterAuth      `json:"masterAuth,omitempty"`
	TokenRefreshTimestamp *metav1.Time     `json:"tokenRefreshTimestamp,omitempty"`
}

type ClusterConditionType string

const (
	ClusterBootstrapAgent ClusterConditionType = "ClusterBootstrapAgent"
	ClusterApprove        ClusterConditionType = "ClusterApprove"
	ClusterRegister       ClusterConditionType = "ClusterRegister"
	ClusterReady          ClusterConditionType = "ClusterReady"
	ClusterDelete         ClusterConditionType = "ClusterDelete"
)

type ConditionStatus string

const (
	Success ConditionStatus = "Success"
	Pending ConditionStatus = "Pending"
)

type ClusterCondition struct {
	Type        ClusterConditionType `json:"type,omitempty"`
	Reason      string               `json:"reason,omitempty"`
	Status      ConditionStatus      `json:"status,omitempty"`
	LastUpdated metav1.Time          `json:"lastUpdated,omitempty"`
}

// TargetStatus defines the observed state of Target
type TargetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions                []ClusterCondition `json:"conditions,omitempty"`
	LastTokenRefreshTimestamp *metav1.Time       `json:"lastTokenRefreshTimestamp,omitempty"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Target is the Schema for the targets API
type Target struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TargetSpec   `json:"spec,omitempty"`
	Status TargetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TargetList contains a list of Target
type TargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Target `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Target{}, &TargetList{})
}

func (t *Target) RefreshToken() bool {
	if t.Spec.TokenRefreshTimestamp.IsZero() {
		return false
	}
	if t.Status.LastTokenRefreshTimestamp.IsZero() {
		return true
	}
	return t.Status.LastTokenRefreshTimestamp.Before(t.Spec.TokenRefreshTimestamp)
}
