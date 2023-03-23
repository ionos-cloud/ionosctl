---
description: Get a Kubernetes NodePool
---

# K8sNodepoolGet

## Usage

```text
ionosctl k8s nodepool get [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific NodePool from an existing Kubernetes Cluster. You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --no-headers           When using text output, don't print headers
  -i, --nodepool-id string   The unique K8s Node Pool Id (required)
  -t, --timeout int          Timeout option for waiting for NodePool to be in ACTIVE state [seconds] (default 600)
  -W, --wait-for-state       Wait for specified NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

