---
description: "Update a Kubernetes NodePool (worker Nodes)"
---

# K8sNodepoolUpdate

## Usage

```text
ionosctl compute k8s nodepool update [flags]
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

Update a node pool: resize it, enable or disable autoscaling, upgrade its
Kubernetes version, change labels/annotations, adjust the maintenance window,
attach LANs, or add reserved public IPs.

Autoscaling: set --min-node-count and --max-node-count together to let Managed
Kubernetes automatically scale the pool between those bounds. Set both to 0 to
disable autoscaling and return to a fixed --node-count.

Kubernetes version: --k8s-version is upgrade-only (downgrades are rejected) and
must stay <= the cluster's version. Upgrade the cluster first if needed.

Maintenance window: --maintenance-day and --maintenance-time define the weekly
window during which IONOS may recycle Nodes for patches/upgrades.

Reserved public IPs (--public-ips): IPs must be reserved in the same location as
the pool's Data Center. Provide one more IP than the pool's maximum node count -
the extra IP is used while Nodes are being rebuilt/rolled.

LANs: --lan-ids adds LANs to the ones already attached (it does not replace
them). To also set a route network and gateway IP on a LAN, use
`ionosctl compute k8s nodepool lan add` per LAN.

Use `--wait` (`-w`) to block until the node pool returns to the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --annotation-key string        Annotation key. Must be set together with --annotation-value (DEPRECATED: Use --labels, --annotations options instead!)
      --annotation-value string      Annotation value. Must be set together with --annotation-key (DEPRECATED: Use --labels, --annotations options instead!)
  -A, --annotations stringToString   Kubernetes annotations for the pool's Nodes. Overwrites any existing annotations. Format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string            The unique K8s Cluster Id (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [NodePoolId Name K8sVersion NodeCount DatacenterId State CpuFamily ServerType StorageType LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId]
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                    Level of detail for response objects (default 1)
      --dhcp                         Whether Nodes obtain an IP via DHCP on the LANs being attached. e.g. --dhcp=true, --dhcp=false (default true)
  -F, --filters strings              Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
      --k8s-version string           New Kubernetes version for the worker Nodes, e.g. 1.29.5. Upgrade-only (downgrades are rejected) and must be <= the cluster's version
      --label-key string             Label key. Must be set together with --label-value (DEPRECATED: Use --labels, --annotations options instead!)
      --label-value string           Label value. Must be set together with --label-key (DEPRECATED: Use --labels, --annotations options instead!)
  -L, --labels stringToString        Kubernetes labels for the pool's Nodes. Overwrites any existing labels. Format: --labels KEY=VALUE,KEY=VALUE (default [])
      --lan-ids ints                 IDs of existing LANs to attach to the worker Nodes. These are added to the LANs already attached, not replacing them. Usage: --lan-ids 1,2
      --limit int                    Maximum number of items to return per request (default 50)
      --maintenance-day string       Day of the week for the maintenance window, in English (e.g. Monday, Saturday). Set together with --maintenance-time
      --maintenance-time string      Start time of the maintenance window in HH:mm:ss (24-hour) format, e.g. 08:00:00. Set together with --maintenance-day
      --max-node-count int           Autoscaling upper bound: maximum number of Nodes the pool may scale out to. Set together with --min-node-count. Set both to 0 to disable autoscaling (default 1)
      --min-node-count int           Autoscaling lower bound: minimum number of Nodes the pool may scale in to. Set together with --max-node-count. Set both to 0 to disable autoscaling (default 1)
      --no-headers                   Don't print table headers when table output is used
      --node-count int               Fixed number of worker Nodes in the pool. Ignored while autoscaling is enabled (min/max node count) (default 1)
  -i, --nodepool-id string           The unique K8s Node Pool Id (required)
      --offset int                   Number of items to skip before starting to collect the results
      --order-by string              Property to order the results by
  -o, --output string                Desired output format [text|json|api-json] (default "text")
      --public-ips strings           Reserved public IPs for the Nodes. IPs must be reserved in the same location as the pool's Data Center. Provide one more IP than the pool's maximum node count (the extra is used while rebuilding Nodes). Usage: --public-ips IP1,IP2
      --query string                 JMESPath query string to filter the output
  -q, --quiet                        Quiet output
      --server-type string           Compute-engine server type for the Nodes: 'DedicatedCore' (dedicated physical CPU cores) or 'VCPU' (shared vCPU cores, typically cheaper). Can be one of: DedicatedCore, VCPU
  -t, --timeout int                  Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                         Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Resize a node pool to 5 Nodes
ionosctl compute k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count 5

# Enable autoscaling between 2 and 8 Nodes and move the maintenance window to Sunday 03:00
ionosctl compute k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID \
  --min-node-count 2 --max-node-count 8 --maintenance-day Sunday --maintenance-time 03:00:00
```

