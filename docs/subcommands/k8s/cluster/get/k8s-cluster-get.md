---
description: Get a Kubernetes Cluster
---

# K8sClusterGet

## Usage

```text
ionosctl k8s cluster get [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -i, --cluster-id string   The unique K8s Cluster Id (required)
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
      --no-headers          When using text output, don't print headers
  -t, --timeout int         Timeout option for waiting for Cluster to be in ACTIVE state [seconds] (default 600)
  -W, --wait-for-state      Wait for specified Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl k8s cluster get --cluster-id CLUSTER_ID
```

