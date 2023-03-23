---
description: Add a Kubernetes NodePool LAN
---

# K8sNodepoolLanAdd

## Usage

```text
ionosctl k8s nodepool lan add [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `add` command:

```text
[a]
```

## Description

Use this command to add a Node Pool LAN into an existing Node Pool.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LanId Dhcp RoutesNetwork RoutesGatewayIp] (default [LanId,Dhcp,RoutesNetwork,RoutesGatewayIp])
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                 Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
      --gateway-ip ip        IPv4 or IPv6 Gateway IP for the route. Must be set with --network flag
  -i, --lan-id int           The unique LAN Id of existing LANs to be attached to worker Nodes (required)
      --network string       IPv4 or IPv6 CIDR to be routed via the interface. Must be set with --gateway-ip flag
      --nodepool-id string   The unique K8s Node Pool Id (required)
```

## Examples

```text
ionosctl k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

