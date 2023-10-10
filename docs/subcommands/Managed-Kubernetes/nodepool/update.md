---
description: "Update a Kubernetes NodePool"
---

# K8sNodepoolUpdate

## Usage

```text
ionosctl k8s nodepool update [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster. You can also add reserved public IP addresses to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. The array must contain one extra IP than maximum number of Nodes could be. The extra provided IP Will be used during rebuilding of Nodes.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Note: If you want to attach multiple LANs to Node Pool, use `--lan-ids=LAN_ID1,LAN_ID2` flag. If you want to also set a Route Network, Route GatewayIp for LAN use `ionosctl k8s nodepool lan add` command for each LAN you want to add.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --annotation-key string        Annotation key. Must be set together with --annotation-value (DEPRECATED: Use --labels, --annotations options instead!)
      --annotation-value string      Annotation value. Must be set together with --annotation-key (DEPRECATED: Use --labels, --annotations options instead!)
  -A, --annotations stringToString   Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host url (default "https://api.ionos.com")
      --cluster-id string            The unique K8s Cluster Id (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32                  Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                         Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
      --k8s-version string           The K8s version for the NodePool. K8s version downgrade is not supported
      --label-key string             Label key. Must be set together with --label-value (DEPRECATED: Use --labels, --annotations options instead!)
      --label-value string           Label value. Must be set together with --label-key (DEPRECATED: Use --labels, --annotations options instead!)
  -L, --labels stringToString        Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
      --lan-ids ints                 Collection of LAN Ids of existing LANs to be attached to worker Nodes. It will be added to the existing LANs attached
      --maintenance-day string       The day of the week for Maintenance Window has the English day format as following: Monday or Saturday
      --maintenance-time string      The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00
      --max-node-count int           The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count (default 1)
      --min-node-count int           The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count (default 1)
      --node-count int               The number of worker Nodes that the NodePool should contain (default 1)
  -i, --nodepool-id string           The unique K8s Node Pool Id (required)
  -o, --output string                Desired output format [text|json|api-json] (default "text")
      --public-ips strings           Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2
  -q, --quiet                        Quiet output
  -t, --timeout int                  Timeout option for waiting for NodePool to be in ACTIVE state [seconds] (default 600)
  -v, --verbose                      Print step-by-step process when running command
  -W, --wait-for-state               Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count NODE_COUNT
```

