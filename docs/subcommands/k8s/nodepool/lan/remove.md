---
description: Remove a Kubernetes NodePool LAN
---

# K8sNodepoolLanRemove

## Usage

```text
ionosctl k8s nodepool lan remove [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `remove` command:

```text
[r]
```

## Description

This command removes a Kubernetes Node Pool LAN from a Node Pool.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
  -a, --all                  Remove all FK8s Nodepool Lans.
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -i, --lan-id int           The unique LAN Id of existing LANs to be detached from worker Nodes (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

