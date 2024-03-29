---
template: main.html
---

# A/B Experiments

A/B testing an application's backend component is challenging.
A/B testing typically relies on business metrics computed by a frontend, user-facing, service.
Metric values often depend on one or more interactions with backend (not user-facing) components.
To A/B test a backend component, it is necessary to be able to associate a metric value (computed by the frontend) to the version of the backend component that contributed to its computation.
The challenge is that the frontend service often does not know which version of the backend component processed a given request.

To address this challenge, Iter8 introduces an A/B/n SDK which provides a frontend service with two APIs:

a. **Lookup()** - identifies a version of a backend component to send a request to

b. **WriteMetric()** - associates a metric with a backend component

This SDK, implemented using gRPC, can be used from a number of frontend implementation languages including *Node.js*, *Python*, *Ruby*, and *Go*, among others. Details of the Iter8 SDK are documented in the [gRPC protoc file](https://github.com/iter8-tools/iter8/blob/v0.13.0/abn/grpc/abn.proto).

This tutorial describes an A/B testing experiment for a backend component.
Example implementations of frontend components are provided in *Node.js* and *Go*.

<p align='center'>
<img alt-text="A/B/n experiment" src="../images/abn.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
 
## Launch Iter8 A/B/n service

Deploy the Iter8 A/B/n service. When deploying the service, specify which Kubernetes resources to watch for each application. To watch for versions of the *backend* application in the *default* namespace, configure the service to watch for service and deployment resources:

```shell
helm install --repo https://iter8-tools.github.io/iter8 iter8-abn abn \
--set "apps.default.backend.resources={service,deployment}"
```

??? warn "Assumptions"
    To simplify specification, Iter8 assumes certain conventions:

    - resources of all versions are deployed to the same namespace
    - there is only one resource of each resource type among the resources of a version
    - all resources that comprise the baseline version are named as: _&lt;application\_name&gt;_
    - all resources that comprise the i<sup>th</sup> candidate version are named as: _&lt;application\_name&gt;-candidate-&lt;i&gt;_

## Deploy the sample application

Deploy both the frontend and backend components of the application as described in each tab:

=== "frontend"
    Install the frontend service using an implementation in the language of your choice:

    === "node"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-node:0.13
        kubectl expose deployment frontend --name=frontend --port=8090
        ```

    === "Go"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-go:0.13
        kubectl expose deployment frontend --name=frontend --port=8090
        ```
    
    The frontend service is implemented to call **Lookup()** before each call to the backend service. It sends its request to the recommended backend service.

=== "backend"
    Deploy version *v1* of the *backend* component as track *backend*.

    ```shell
    kubectl create deployment backend --image=iter8/abn-sample-backend:0.13-v1
    kubectl label deployment backend app.kubernetes.io/version=v1

    kubectl expose deployment backend --name=backend --port=8091
    ```

Before calling the backend, the frontend uses *Lookup()* to identify the track to send requests to. Since there is only one version of the backend deployed, all requests will be sent to it.

## Generate load

Generate load. In separate shells, port-forward requests to the frontend service and generate load for multiple users.  For example:
    ```shell
    kubectl port-forward service/frontend 8090:8090
    ```
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/main/samples/abn-sample/generate_load.sh | sh -s --
    ```

Note that the the names `foo` and `foobar` are examples. They may be mapped to the same track label -- since we are using 

## Deploy a candidate version

Deploy version *v2* of the *backend* component as track *backend-candidate-1*.

```shell
kubectl create deployment backend-candidate-1 --image=iter8/abn-sample-backend:0.13-v2
kubectl label deployment backend-candidate-1 app.kubernetes.io/version=v2

kubectl expose deployment backend-candidate-1 --name=backend-candidate-1 --port=8091
```

Until the candidate version is ready; that is, until all expected resources are deployed and available, calls to *Lookup()* will continue to return only the *backend* track.
Once the candidate version is ready, *Lookup()* will return both tracks so that requests will be distributed between them.

## Launch experiment

```shell
iter8 k launch \
--set abnmetrics.application=default/backend \
--set "tasks={abnmetrics}" \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

??? note "About this experiment"
    This experiment periodically (in this case, once a minute) reads the `abn` metrics associated with the *backend* application component in the *default* namespace. These metrics are written by the frontend service using the *WriteMetric()* interface as a part of processing user requests.

## Inspect experiment report

Inspect the metrics:

```shell
iter8 k report
```

??? note "Sample output from report"
    ```
    Experiment summary:
    *******************

    Experiment completed: true
    No task failures: true
    Total number of tasks: 1
    Number of completed tasks: 1
    Number of completed loops: 3

    Latest observed values for metrics:
    ***********************************

    Metric                   | backend (v1) | backend-candidate-1 (v2)
    -------                  | -----        | -----
    abn/sample_metric/count  | 35.00        | 28.00
    abn/sample_metric/max    | 99.00        | 100.00
    abn/sample_metric/mean   | 56.31        | 52.79
    abn/sample_metric/min    | 0.00         | 1.00
    abn/sample_metric/stddev | 28.52        | 31.91
    ```
The output allows you to compare the versions against each other and select a winner. Since the experiment runs periodically, you should expect the values in the report to change over time.

Once a winner is identified, the experiment can be terminated and the winner can be promoted and the candidate versions can be deleted.

To delete the experiment:

```shell
iter8 k delete
```

## Promote candidate version

Delete the candidate version:

```shell
kubectl delete deployment backend-candidate-1 
kubectl delete service backend-candidate-1
```

Update the version of the baseline track:

```shell
kubectl set image deployment/backend abn-sample-backend=iter8/abn-sample-backend:0.13-v2
kubectl label --overwrite deployment/backend app.kubernetes.io/version=v2

# no change in service
# kubectl expose deployment backend --name=backend --port=8091
```

## Cleanup

### Delete sample application

```shell
kubectl delete \
deploy/frontend deploy/backend deploy/backend-candidate-1 \
service/frontend service/backend service/backend-candidate-1
```

### Uninstall the A/B/n service

```shell
helm delete iter8-abn
```
