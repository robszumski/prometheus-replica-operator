package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PrometheusReplicaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PrometheusReplica `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PrometheusReplica struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PrometheusReplicaSpec   `json:"spec"`
	Status            PrometheusReplicaStatus `json:"status,omitempty"`
}

type PrometheusReplicaSpec struct {
	// Fill me
}
type PrometheusReplicaStatus struct {
	// Fill me
}
