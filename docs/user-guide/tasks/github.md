---
template: main.html
---

# github

Trigger GitHub workflows via a [repository_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch).

A `repository_dispatch` will trigger workflows in the default branch of the GitHub repository. By default, the experiment report will also be sent.

## Usage Example

```shell
iter8 k launch \
--set "tasks={http,assess,github}" \
--set http.url=http://httpbin.default/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set github.owner=<GitHub owner> \
--set github.repo=<GitHub repository> \
--set github.token=<GitHub token> \
--set runner=job
```

See [here](../../tutorials/integrations/ghactions.md#use-iter8-to-trigger-a-github-actions-workflow) for a more in-depth tutorial.

## Parameters

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| owner | string | Yes | N/A | Owner of the GitHub repository |
| repo | string | Yes | N/A | GitHub repository |
| token | string | Yes | N/A | Authorization token |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/iter8/iter8-0.13.7/templates/notify/_payload-github.tpl](https://raw.githubusercontent.com/iter8-tools/iter8/iter8-0.13.7/templates/notify/_payload-github.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the task and experiment should not fail if the task cannot successfully send the request |
| if | string | No | N/A | An if condition that can be control when the task is run in a [multi-looped experiment](../../getting-started/concepts.md#runner). To learn more, see [here](#if-parameter). |

## Default payload

A `repository_dispatch` requires a payload that contains the type of the event. 

The [default payload template](https://raw.githubusercontent.com/iter8-tools/iter8/iter8-0.13.7/templates/notify/_payload-github.tpl) will set the `event_type` to `iter8`. In addition, it will also provide the experiment report in the `client_payload`, which means that this data will be accessible in the GitHub workflow via `${{ toJson(github.event.client_payload) }}`.

However, if you would like to use a different payload template, simply set a `payloadTemplateURL` and Iter8 will not use the default.

## `if` parameter

The `if` parameter is used to control when the task is run in a [multi-looped experiment](../../getting-started/concepts.md#runner). For example, if you would like for the `github` task to only run at the 10th loop instead of every loop, you can do the following:

```diff
  iter8 k launch \
  --set "tasks={http,assess,github}" \
  --set http.url=http://httpbin.default/get \
  --set assess.SLOs.upper.http/latency-mean=50 \
  --set assess.SLOs.upper.http/error-count=0 \
  --set github.owner=<GitHub owner> \
  --set github.repo=<GitHub repository> \
  --set github.token=<GitHub token> \
  --set github.if="Result.NumLoops == 10"
  --set runner=job
```

You may use any field in the `Result` object for your logic.