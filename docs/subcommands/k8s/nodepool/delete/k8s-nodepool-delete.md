---
description: Delete a Kubernetes NodePool
---

# K8sNodepoolDelete

## Usage

```text
ionosctl k8s nodepool delete [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -a, --all                  Delete all the Kubernetes Node Pools within an existing Kubernetes Nodepools.
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -i, --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

