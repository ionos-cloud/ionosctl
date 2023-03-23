---
description: Get a Kubernetes Node
---

# K8sNodeGet

## Usage

```text
ionosctl k8s node get [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Kubernetes Node.You can wait for the Node to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --no-headers           When using text output, don't print headers
  -i, --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -t, --timeout int          Timeout option for waiting for Node to be in ACTIVE state [seconds] (default 600)
  -W, --wait-for-state       Wait for specified Node to be in ACTIVE state
```

## Examples

```text
ionosctl k8s node get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

