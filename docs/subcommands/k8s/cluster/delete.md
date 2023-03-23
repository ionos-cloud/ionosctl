---
description: Delete a Kubernetes Cluster
---

# K8sClusterDelete

## Usage

```text
ionosctl k8s cluster delete [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

You can wait for Request for the Cluster deletion to be executed using `--wait-for-request` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -a, --all                 Delete all the Kubernetes clusters.
  -i, --cluster-id string   The unique K8s Cluster Id (required)
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int         Timeout option for waiting for Request [seconds] (default 600)
  -w, --wait-for-request    Wait for the Request for Cluster deletion to be executed
```

## Examples

```text
ionosctl k8s cluster delete --cluster-id CLUSTER_ID
```

