---
template: main.html
---

# Metrics

## Provider

A **provider** in Iter8 is a data source that supplies metric values.

## Fully qualified names

Metrics are scoped by providers. Providers have unique names, and within the scope of a provider, metrics have unique names. The fully qualified name of a metric refers to the string of the form `<provider name>/<metric name>`.

Following are some examples of fully qualified metric names that appear in Iter8 tutorials.

1. `http/latency-mean`
2. `grpc/error-rate`
3. `istio-prom/latency-mean`

## Built-in metrics provider

Iter8 has built-in metrics providers, namely, [`http`](../tasks/http.md#metrics) and [`grpc`](../tasks/grpc.md#metrics).

## Custom metrics provider

You can use metrics from any (RESTful) database in Iter8 experiments. Metrics fetched by Iter8 from databases are also referred to as custom metrics. See [here](../tasks/custommetrics.md) to learn more about custom metrics.

## Metric types

Iter8 defines `counter` and `gauge` metric types which are analogous to the corresponding [metric types defined by Prometheus](https://prometheus.io/docs/concepts/metric_types/). We quote from the Prometheus documentation below for their definitions.

???+ note "Counter"
    A counter is a cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart. 
    
    For example, you can use a counter to represent the number of requests served, tasks completed, or errors. Do not use a counter to expose a value that can decrease. For example, do not use a counter for the number of currently running processes; instead use a gauge.

???+ note "Gauge"

    A gauge is a metric that represents a single numerical value that can arbitrarily go up and down. 
    
    Gauges are typically used for measured values like temperatures or current memory usage, but also "counts" that can go up and down, like the number of concurrent requests.

