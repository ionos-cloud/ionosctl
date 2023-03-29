---
description: Recreate a Kubernetes Node
---

# K8sNodeRecreate

## Usage

```text
ionosctl k8s node recreate [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `recreate` command:

```text
[r]
```

## Description

You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -i, --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s node recreate --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

