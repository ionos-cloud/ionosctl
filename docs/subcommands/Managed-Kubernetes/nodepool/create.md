---
description: "Create a Kubernetes NodePool (worker Nodes)"
---

# K8sNodepoolCreate

## Usage

```text
ionosctl compute k8s nodepool create [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `create` command:

```text
[c]
```

## Description

Create a node pool of worker Nodes inside an existing Managed Kubernetes cluster.
The cluster must be in state ACTIVE first. Every Node in the pool is provisioned
into the Data Center named by --datacenter-id and shares the same hardware shape
(--cores, --ram, --storage-size, --storage-type, --cpu-family) and Kubernetes
version.

Name: up to 63 characters, must begin and end with an alphanumeric character,
with dashes, underscores, dots and alphanumerics in between.

Kubernetes version: if --k8s-version is not set, the parent cluster's version is
used. The pool version must be less than or equal to the cluster version.

CPU family: if --cpu-family is omitted, the first CPU family available in the
Data Center's location is chosen automatically. --server-type (VCPU or
DedicatedCore) selects the compute engine server type for the Nodes.

Autoscaling and reserved public IPs are not configurable at create time via
dedicated flags; set them afterwards with `ionosctl compute k8s nodepool update`,
or pass a full JSON body with --json-properties.

Networking: attach existing LANs with --lan-ids=LAN_ID1,LAN_ID2. To also set a
route network and gateway IP on a LAN, use `ionosctl compute k8s nodepool lan add`
per LAN after creation.

Use `--wait` (`-w`) to block until the node pool reaches the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* Datacenter Id

## Options

```text
  -A, --annotations stringToString   Kubernetes annotations propagated onto the pool's Nodes. Overwrites any existing annotations. Format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string     Compute availability zone for the Nodes. AUTO lets IONOS place them; ZONE_1 / ZONE_2 pin them to a specific zone (default "AUTO")
      --cluster-id string            The unique K8s Cluster Id (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [NodePoolId Name K8sVersion NodeCount DatacenterId State CpuFamily ServerType StorageType LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId]
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                    Number of CPU cores per Node (default 2)
      --cpu-family string            CPU family for the Nodes (e.g. INTEL_SKYLAKE, INTEL_XEON), constrained by the Data Center's location. If not set, the first CPU family available in that location (as returned by the API) is used
      --datacenter-id string         The unique Data Center Id (required)
  -D, --depth int                    Level of detail for response objects (default 1)
      --dhcp                         Whether Nodes obtain an IP on the attached LANs via DHCP. Applies to the LANs given in --lan-ids. e.g. --dhcp=true, --dhcp=false (default true)
  -F, --filters strings              Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
      --json-properties string       Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example      If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --k8s-version string           Kubernetes version for the worker Nodes, e.g. 1.29.5. Must be <= the cluster's version. If not set, the cluster's version is used
  -L, --labels stringToString        Kubernetes labels propagated onto the pool's Nodes (usable for scheduling/nodeSelectors). Overwrites any existing labels. Format: --labels KEY=VALUE,KEY=VALUE (default [])
      --lan-ids nodepool lan add     IDs of existing LANs (in the same Data Center) to attach to the worker Nodes, e.g. --lan-ids 1,2. Use nodepool lan add to also set routes on a LAN
      --limit int                    Maximum number of items to return per request (default 50)
  -n, --name string                  Name of the node pool. Up to 63 characters; must start and end with an alphanumeric character, dashes, underscores, dots and alphanumerics in between (default "UnnamedNodePool")
      --no-headers                   Don't print table headers when table output is used
      --node-count int               Number of worker Nodes in the pool. Minimum 1; the maximum depends on your contract and resource availability (default 1)
      --offset int                   Number of items to skip before starting to collect the results
      --order-by string              Property to order the results by
  -o, --output string                Desired output format [text|json|api-json] (default "text")
      --query string                 JMESPath query string to filter the output
  -q, --quiet                        Quiet output
      --ram string                   RAM per Node. Minimum 2048MB and must be a multiple of 1024MB. Accepts a unit suffix, e.g. --ram 2048, --ram 2048MB or --ram 4GB (default "2048")
      --server-type string           Compute-engine server type for the Nodes: 'DedicatedCore' (dedicated physical CPU cores) or 'VCPU' (shared vCPU cores, typically cheaper). Can be one of: DedicatedCore, VCPU
      --storage-size string          Per-Node boot storage size in GB. Accepts a unit suffix, e.g. --storage-size 10 or --storage-size 10GB. The maximum is bounded by your contract limit (default "10")
      --storage-type string          Type of the per-Node boot storage: HDD (default) or SSD (default "HDD")
  -t, --timeout int                  Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                         Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a minimal node pool (1 Node, defaults for hardware and version) and wait for it
ionosctl compute k8s nodepool create --cluster-id CLUSTER_ID --datacenter-id DATACENTER_ID --name pool-a --wait

# Create a 3-Node pool with a specific shape, SSD storage, pinned Kubernetes version and two attached LANs
ionosctl compute k8s nodepool create --cluster-id CLUSTER_ID --datacenter-id DATACENTER_ID --name workers \
  --node-count 3 --cores 4 --ram 8GB --storage-type SSD --storage-size 100 --k8s-version 1.29.5 --lan-ids 1,2
```

