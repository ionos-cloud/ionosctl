---
description: "Create a Kubernetes NodePool"
---

# K8sNodepoolCreate

## Usage

```text
ionosctl k8s nodepool create [flags]
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

Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character with dashes, underscores, dots, and alphanumerics between. Regarding the Kubernetes Version for the NodePool, if not set via flag, it will be used the default one: `ionosctl k8s version get`.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Note: If you want to attach multiple LANs to Node Pool, use `--lan-ids=LAN_ID1,LAN_ID2` flag. If you want to also set a Route Network, Route GatewayIp for LAN use `ionosctl k8s nodepool lan add` command for each LAN you want to add.

Required values to run a command (for Public Kubernetes Cluster):

* K8s Cluster Id
* Datacenter Id

Required values to run a command (for Private Kubernetes Cluster):

* K8s Cluster Id
* Datacenter Id

## Options

```text
  -A, --annotations stringToString   Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string     The compute Availability Zone in which the Node should exist (default "AUTO")
      --cluster-id string            The unique K8s Cluster Id (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily ServerType StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                    The total number of cores for the Node (default 2)
      --cpu-family string            CPU Type. If the flag is not set, the CPU Family will be chosen based on the location of the Datacenter. It will always be the first CPU Family available, as returned by the API
      --datacenter-id string         The unique Data Center Id (required)
  -D, --depth int32                  Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                         Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
      --json-properties string       Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example      If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --k8s-version string           The K8s version for the NodePool. If not set, the default one will be used
  -L, --labels stringToString        Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
      --lan-ids ints                 Collection of LAN Ids of existing LANs to be attached to worker Nodes
  -n, --name string                  The name for the K8s NodePool (default "UnnamedNodePool")
      --no-headers                   Don't print table headers when table output is used
      --node-count int               The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability (default 1)
  -o, --output string                Desired output format [text|json|api-json] (default "text")
  -q, --quiet                        Quiet output
      --ram string                   RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB (default "2048")
      --server-type string           The type of server for the Kubernetes node pool can be either'DedicatedCore' (nodes with dedicated CPU cores) or 'VCPU' (nodes with shared CPU cores).This selection corresponds to the server type for the compute engine.. Can be one of: DedicatedCore, VCPU
      --storage-size string          The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit (default "10")
      --storage-type string          Storage Type (default "HDD")
  -t, --timeout int                  Timeout option for waiting for NodePool to be in ACTIVE state[seconds] (default 600)
  -v, --verbose                      Print step-by-step process when running command
  -W, --wait-for-state               Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool create --datacenter-id DATACENTER_ID --cluster-id CLUSTER_ID --name NAME
```

