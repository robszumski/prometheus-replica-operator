# Prometheus Replica Operator

A Kubernetes Operator for Prometheus + Thanos, build on top of the [Operator SDK}().

[This blog post shows how I built](https://robszumski.com/building-an-operator/) the Prometheus Replica Operator, which is the first complex Go program I’ve written. While the Go code itself may not be the best example of software engineering, I wanted to try out my team’s Operator SDK. As the Product Manager for the Operator SDK, I want to have first hand knowledge of our tool. As an early employee of CoreOS, I do have a lot of knowledge about Operators, and I want to pass on some best practices through the post.

## What the PRO does

The PRO will install and configure a full monitoring stack on a Kubernetes cluster, using [Prometheus](https://prometheus.io/) for ingesting, storing and querying the time series data. Archival data is automatically sent to a cloud storage bucket with [Thanos(https://github.com/improbable-eng/thanos).

## Example

The Operator watches for `PrometheusReplica` objects, such as this one:

```
apiVersion: "prometheus.robszumski.com/v1alpha1"
kind: "PrometheusReplica"
metadata:
  name: "example"
spec:
  configMap: prometheus-config
  highlyAvailable: true
  baseDomain: "ingress.example.com"
  metrics:
    retention: 24h
    blockDuration: 1h
  bucketSecret: s3-bucket
```
  
And configures the entire monitoring stack:

```
INFO[0000] Go Version: go1.10.2
INFO[0000] Go OS/Arch: darwin/amd64
INFO[0000] operator-sdk Version: 0.0.5+git
INFO[0000] Watching prometheus.robszumski.com/v1alpha1, PrometheusReplica, default, 5
INFO[0000] starting prometheusreplicas controller
...detected object...
INFO[0000] Parsing PrometheusReplica example in default
INFO[0000] Updating PrometheusReplica status for example
INFO[0000] Status of PrometheusReplica example is now Install
INFO[0000] Creating Prometheus StatefulSet for example
INFO[0000]   StatefulSet: Translating HighlyAvailable to 2 replicas
INFO[0000]   StatefulSet: Setting overall metrics retention to 24h
INFO[0000]   StatefulSet: Setting duration until upload to storage bucket to 1h
INFO[0000]   StatefulSet: Using Prometheus config from ConfigMap prometheus-config
INFO[0000]   StatefulSet: Using bucket parameters from Secret s3-bucket
INFO[0001] Creating Prometheus service for example
INFO[0001] Creating Thanos peers service for example
INFO[0001] Creating Thanos store StatefulSet for example
INFO[0001]   StatefulSet: Using bucket parameters from Secret s3-bucket
INFO[0001] Creating Thanos store service for example
INFO[0001] Creating Thanos query Deployment for example
INFO[0001]   Deployment: Using bucket parameters from Secret s3-bucket
INFO[0001]   Deployment: Translating HighlyAvailable to 2 replicas
INFO[0001] Creating Thanos query service for example
INFO[0001] Checking desired vs actual state for components of PrometheusReplica example
INFO[0002] Creating Prometheus StatefulSet for example
INFO[0002]   StatefulSet: Translating HighlyAvailable to 2 replicas
INFO[0002]   StatefulSet: Setting overall metrics retention to 24h
INFO[0002]   StatefulSet: Setting duration until upload to storage bucket to 1h
INFO[0002]   StatefulSet: Using Prometheus config from ConfigMap prometheus-config
INFO[0002]   StatefulSet: Using bucket parameters from Secret s3-bucket
INFO[0002]   Checking StatefulSet for Prometheus
INFO[0004] Parsing PrometheusReplica example in default
INFO[0004] Updating PrometheusReplica status for example
INFO[0004] Status of PrometheusReplica example is now Creating
INFO[0005] Checking desired vs actual state for components of PrometheusReplica example
INFO[0005]   Checking StatefulSet for Prometheus
INFO[0005]   Checking Deployment for Thanos query
...create is now done...
INFO[0005] Parsing PrometheusReplica example in default
INFO[0005] Updating PrometheusReplica status for example
INFO[0005] Status of PrometheusReplica example is now Running
INFO[0006] Checking desired vs actual state for components of PrometheusReplica example
INFO[0006]   Checking StatefulSet for Prometheus
INFO[0006]   Checking Deployment for Thanos query
...looping...
INFO[0006] Parsing PrometheusReplica example in default
INFO[0006] Updating PrometheusReplica status for example
INFO[0007] Checking desired vs actual state for components of PrometheusReplica example
INFO[0007]   Checking StatefulSet for Prometheus
INFO[0008]   Checking Deployment for Thanos query
...loop forever...
  ```
  
## Install

First, install the CRD:

```
$ kubectl create -f http://todo:insert Github raw URL
```

Then run the Operator:

```
$ kubectl create -f http://todo:insert Github raw URL
```

Last, create the `PrometheusReplica` object:

```
$ kubectl create -f http://todo:insert Github raw URL
```

If everything worked correctly, you should see the Pods, Services, Deployments and StatefulSets created. The status of the `PrometheusReplica` should also have this information:

```
$ kubectl get prometheusreplica/todo -o yaml
```
