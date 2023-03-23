---
description: List Kubernetes NodePool LANs
---

# K8sNodepoolLanList

## Usage

```text
ionosctl k8s nodepool lan list [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of all contained NodePool LANs in a selected Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LanId Dhcp RoutesNetwork RoutesGatewayIp] (default [LanId,Dhcp,RoutesNetwork,RoutesGatewayIp])
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -M, --max-results int32    The maximum number of elements to return
      --no-headers           When using text output, don't print headers
      --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s nodepool lan list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

