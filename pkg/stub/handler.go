package stub

import (
	"fmt"
	//"reflect"
	"context"

	v1alpha1 "github.com/robszumski/prometheus-replica-operator/pkg/apis/prometheus/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/labels"
	"github.com/sirupsen/logrus"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.PrometheusReplica:
		PrometheusReplica := o

		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			return nil
		}

		// SANITY CHECK
		logrus.Infof("Parsing PrometheusReplica %s in %s", PrometheusReplica.Name, PrometheusReplica.Namespace)


        // INSTALL
		// Create the Prometheus StatefulSet if it doesn't exist
		ssProm := statefulSetForPrometheus(PrometheusReplica)
		err := sdk.Create(ssProm)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create Prometheus StatefulSet: %v", err)
		}

		// // Create the Prometheus ConfigMap if it doesn't exist
		// cmProm := configmapSetForPrometheus(PrometheusReplica)
		// err := sdk.Create(cmProm)
		// if err != nil && !apierrors.IsAlreadyExists(err) {
		// 	return fmt.Errorf("failed to create Prometheus ConfigMap: %v", err)
		// }

		// UPDATE
		// Ensure the deployment size is the same as the spec
		// err = sdk.Get(dep)
		// if err != nil {
		// 	return fmt.Errorf("failed to get deployment: %v", err)
		// }
		// size := PrometheusReplica.Spec.Size
		// if *dep.Spec.Replicas != size {
		// 	dep.Spec.Replicas = &size
		// 	err = sdk.Update(dep)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to update deployment: %v", err)
		// 	}
		// }

		// Update the PrometheusReplica status with the pod names
		// podList := podList()
		// labelSelector := labels.SelectorFromSet(labelsForPrometheusReplica(PrometheusReplica.Name)).String()
		// listOps := &metav1.ListOptions{LabelSelector: labelSelector}
		// err = sdk.List(PrometheusReplica.Namespace, podList, sdk.WithListOptions(listOps))
		// if err != nil {
		// 	return fmt.Errorf("failed to list pods: %v", err)
		// }
		// podNames := getPodNames(podList.Items)
		// if !reflect.DeepEqual(podNames, PrometheusReplica.Status.Nodes) {
		// 	PrometheusReplica.Status.Nodes = podNames
		// 	err := sdk.Update(PrometheusReplica)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to update PrometheusReplica status: %v", err)
		// 	}
		// }
	}
	return nil
}

// statefulSetForPrometheus returns a PrometheusReplica StatefulSet object
func statefulSetForPrometheus(pr *v1alpha1.PrometheusReplica) *appsv1.StatefulSet {
	ls := labelsForPrometheusReplica(pr.Name)
	retention := pr.Spec.Metrics.Retention
	blockDuration := pr.Spec.Metrics.BlockDuration

	logrus.Infof("Creating Prometheus StatefulSet for %s", pr.Name)

	var replicas int32
	if pr.Spec.HighlyAvailable {
	    replicas = 2
	    logrus.Infof("StatefulSet: Translating HighlyAvailable to %d replicas", replicas)
	} else {
		replicas = 1
		logrus.Infof("StatefulSet: No HA. Starting %t replica", replicas)
	}

	logrus.Infof("StatefulSet: Setting overall metrics retention to %s", retention)
	logrus.Infof("StatefulSet: Setting duration until upload to storage bucket to %s", blockDuration)

	dep := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pr.Name,
			Namespace: pr.Namespace,
			Labels:    ls,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: fmt.Sprintf("%s-prometheus", pr.Name),
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
					Annotations: map[string]string{"prometheus.io/scrape": "true","prometheus.io/port": "10902"},
				},
				Spec: v1.PodSpec{
					Affinity: &v1.Affinity{
						PodAntiAffinity: &v1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
								{
									LabelSelector: &metav1.LabelSelector{
										MatchExpressions: []metav1.LabelSelectorRequirement{
											{
												Key:      "app",
												Operator: metav1.LabelSelectorOpIn,
												Values:   []string{"prometheus"},
											},
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
					Containers: []v1.Container{{
						Image:   "quay.io/prometheus/prometheus:v2.0.0",
						Name:    "prometheus",
						Args: []string{
							fmt.Sprintf("--storage.tsdb.retention=%s", retention),
					        "--config.file=/etc/prometheus-shared/prometheus.yml",
					        "--storage.tsdb.path=/var/prometheus",
					        fmt.Sprintf("--storage.tsdb.min-block-duration=%s", blockDuration),
					        fmt.Sprintf("--storage.tsdb.max-block-duration=%s", blockDuration),
					        "--web.enable-lifecycle",
						},
						Ports: []v1.ContainerPort{{
							ContainerPort: 9090,
							Name:          "http",
						}},
						VolumeMounts: []v1.VolumeMount{{
							MountPath: "/etc/prometheus-shared", Name: "config-shared",
						}, {
							MountPath: "/var/prometheus", Name: "data",
						}},
					},{
						Image:   "improbable/thanos:master",
						Name:    "thanos-sidecar",
						Args: []string{
							"sidecar",
					        "--log.level=debug",
					        "--tsdb.path=/var/prometheus",
					        "--prometheus.url=http://127.0.0.1:9090",
					        "--cluster.peers=thanos-peers.default.svc.cluster.local:10900",
					        "--reloader.config-file=/etc/prometheus/prometheus.yml.tmpl",
					        "--reloader.config-envsubst-file=/etc/prometheus-shared/prometheus.yml",
					        "--s3.signature-version2",
						},
						Env: []v1.EnvVar{{
							Name: "S3_BUCKET",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{Name: "s3-bucket"},
									Key:  "s3_bucket",
								},
							},
						},{
							Name: "S3_ENDPOINT",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{Name: "s3-bucket"},
									Key:  "s3_endpoint",
								},
							},
						},{
							Name: "S3_ACCESS_KEY",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{Name: "s3-bucket"},
									Key:  "s3_access_key",
								},
							},
						},{
							Name: "S3_SECRET_KEY",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{Name: "s3-bucket"},
									Key:  "s3_secret_key",
								},
							},
						}},
						Ports: []v1.ContainerPort{{
							ContainerPort: 10902,
							Name:          "http",
						},{
							ContainerPort: 10901,
							Name:          "grpc",
						},{
							ContainerPort: 10900,
							Name:          "cluster",
						}},
						VolumeMounts: []v1.VolumeMount{{
							MountPath: "/etc/prometheus", Name: "config",
						},{
							MountPath: "/etc/prometheus-shared", Name: "config-shared",
						}, {
							MountPath: "/var/prometheus", Name: "data",
						}},
					}},
					Volumes: []v1.Volume{
						{
							Name: "config-shared",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "data",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "config",
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: "prometheus-config",
								},
							},
						},
					},
				},
			},
		},
	}
	addOwnerRefToObject(dep, asOwner(pr))
	return dep
}

// labelsForPrometheusReplica returns the labels for selecting the resources
// belonging to the given PrometheusReplica CR name.
func labelsForPrometheusReplica(name string) map[string]string {
	return map[string]string{"app": "prometheus", "PrometheusReplica_cr": name, "thanos-peer": "true"}
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the PrometheusReplica CR
func asOwner(m *v1alpha1.PrometheusReplica) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		UID:        m.UID,
		Controller: &trueVar,
	}
}

// podList returns a v1.PodList object
func podList() *v1.PodList {
	return &v1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}