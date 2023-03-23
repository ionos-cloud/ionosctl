---
description: Get the kubeconfig file for a Kubernetes Cluster
---

# K8sKubeconfigGet

## Usage

```text
ionosctl k8s kubeconfig get [flags]
```

## Aliases

For `kubeconfig` command:

```text
[cfg config]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the kubeconfig file for a given Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
      --cluster-id string   The unique K8s Cluster Id (required)
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl k8s kubeconfig get --cluster-id CLUSTER_ID
```

