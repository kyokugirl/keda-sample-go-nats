# NATS Jetstream consumer and publisher

A simple docker container that will receive messages from a NATS stream and scale via KEDA.
The receiver will receive a single message at a time (per instance), and sleep for 1 second to simulate performing work.
When adding a massive amount of stream messages, KEDA will drive the container to scale out according to the event source (NATS).

## Pre-requisites

* Kubernetes cluster
* [KEDA 2.0 installed](https://keda.sh/docs/deploy/) on the cluster
* [NATS CLI](https://docs.nats.io/natscli/nats_cli/) (for stream creation - optional if you create stream via other means)

## Setup

This setup will guide you through creating a NATS JetStream stream on the cluster and deploying a consumer with a `ScaledObject` to scale via KEDA.
If you already have NATS Jetstream, you can use your existing streams.

First you should clone the project:

```cli
git clone https://github.com/kyokugirl/keda-go-NATS
cd keda-go-NATS
```

### Creating a NATS stream

#### [Install Helm](https://helm.sh/docs/using_helm/)

#### Install NATS via Helm

Add the NATS Helm repository:

```cli
helm repo add nats https://nats-io.github.io/k8s/helm/
helm repo update
```

##### Helm 3

```cli
helm install nats nats/nats --wait
```

##### Helm 2


```cli
helm install --name nats nats/nats --wait
```

#### Wait for NATS to Deploy

⚠️ Be sure to wait until the deployment has completed before continuing. ⚠️

```cli
kubectl get po

NAME                     READY   STATUS    RESTARTS   AGE
nats-0                   1/1     Running   0          2m
nats-exporter-6b4cf5999c-jqpdl  1/1     Running   0          2m
```

### Deploying a NATS consumer

#### Deploy a consumer

```cli
kubectl apply -f deploy/deploy-consumer.yaml
```

#### Validate the consumer has deployed

```cli
kubectl get deploy
```

You should see `nats-consumer` deployment with 0 pods as there currently aren't any queue messages and for that reason it is scaled to zero.

```cli
NAME                DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nats-consumer       0         0         0            0           3s
```

[This consumer](https://github.com/kyokugirl/keda-sample-go-nats/blob/main/cmd/receive/receive.go) is set to consume one message per instance, sleep for 1 second, and then acknowledge completion of the message.  This is used to simulate work.  The [`ScaledObject` included in the above deployment](deploy/deploy-consumer.yaml) is set to scale to a minimum of 0 replicas on no events, and up to a maximum of 30 replicas on heavy events (optimizing for a queue length of 5 message per replica).  After 30 seconds of no events the replicas will be scaled down (cooldown period).  These settings can be changed on the `ScaledObject` as needed.

### Publishing messages to the stream

#### Deploy the publisher job

The following job will publish 300 messages to the "hello" queue the deployment is listening to. As the queue builds up, KEDA will help the horizontal pod autoscaler add more and more pods until the queue is drained after about 2 minutes and up to 30 concurrent pods.  You can modify the exact number of published messages in the `deploy-publisher-job.yaml` file.

```cli
kubectl apply -f deploy/deploy-publisher-job.yaml
```

#### Validate the deployment scales

```cli
kubectl get deploy -w
```

You can watch the pods spin up and start to process queue messages.  As the message length continues to increase, more pods will be pro-actively added.

You can see the number of messages vs the target per pod as well:

```cli
kubectl get hpa
```

After the queue is empty and the specified cooldown period (a property of the `ScaledObject`, default of 300 seconds) the last replica will scale back down to zero.

## Cleanup resources

```cli
kubectl delete job nats-publish
kubectl delete ScaledObject nats-consumer
kubectl delete deploy nats-consumer
helm delete nats
```
