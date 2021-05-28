---
description: Create a Kubernetes NodePool
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

Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. Regarding the Kubernetes Version for the NodePool, if not set via flag, it will be used the default one: `ionosctl k8s version get`.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id
* Name

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   The compute Availability Zone in which the Node should exist (default "AUTO")
      --cluster-id string          The unique K8s Cluster Id (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps PublicIps AvailableUpgradeVersions] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  The total number of cores for the Node (default 2)
      --cpu-family string          CPU Type (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -f, --force                      Force command to execute without user input
  -h, --help                       help for create
      --k8s-version string         The K8s version for the NodePool. If not set, it will be used the default one
  -n, --name string                The name for the K8s NodePool (required)
      --node-count int             The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability (default 1)
  -o, --output string              Desired output format [text|json] (default "text")
  -q, --quiet                      Quiet output
      --ram string                 RAM size for node, minimum size is 2048MB. Ram size must be set to multiple of 1024MB. e.g. --ram 2048 or --ram 2048MB (default "2048")
      --storage-size int           The total allocated storage capacity of a Node (default 10)
      --storage-type string        Storage Type (default "HDD")
  -t, --timeout int                Timeout option for waiting for NodePool/Request [seconds] (default 600)
  -W, --wait-for-state             Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool create --datacenter-id DATACENTER_ID --cluster-id CLUSTER_ID --name NAME
```

