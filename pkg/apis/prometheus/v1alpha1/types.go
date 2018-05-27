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


// Spec

type PrometheusReplicaSpec struct {
	// Size is the size of the memcached deployment
	Size string          `json:"configMap"`
	HighlyAvailable bool `json:"highlyAvailable"`
	BaseDomain string    `json:"baseDomain"`
	BucketSecret string  `json:"bucketSecret"`
	Metrics string       PrometheusMetricsSpec `json:"metrics"`
}

type PrometheusMetricSpec struct {
	// Size is the size of the memcached deployment
	Retention string   `json:"retention"`
	BlockDuration bool `json:"blockDuration"`
}

// Status

type PrometheusReplicaStatus struct {
	Phase string  `json:"phase"`
	Ouput         PrometheusReplicaStatus `json:"output,omitempty"`
	Local         PrometheusReplicaStatus `json:"local,omitempty"`
}

type PrometheusOutputStatus struct {
	Grafana string `json:"grafana"`
	Query string   `json:"query"`
}

type PrometheusLocalStatus struct {
	Stores string       `json:"grafana"`
	Prometheuses string `json:"prometheuses"`
	Queries string      `json:"queries"`
}

// spec:
// 	config: ...
// 	ha: true
// 	baseDomain: ...
// 	metrics:
// 		retention: 24h
// 		blockDuration: 2m
// 	bucketSecret: ...
// status:
//   phase: running, creating, deleting
// 	output:
// 		grafana: g.e.com or svc.local
// 		query: q.e.com or svc.local
// 	local:
// 		stores:
// 			...
// 		prometheuses:
// 			...
// 		queries:
// 		  ...
