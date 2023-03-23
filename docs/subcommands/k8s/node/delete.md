---
description: Delete a Kubernetes Node
---

# K8sNodeDelete

## Usage

```text
ionosctl k8s node delete [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes Node within an existing Kubernetes NodePool in a Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
  -a, --all                  Delete all the Kubernetes Nodes within an existing Kubernetes NodePool in a Cluster.
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s node delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

