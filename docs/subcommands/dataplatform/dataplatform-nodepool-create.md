---
description: Create a Data Platform NodePool
---

# DataplatformNodepoolCreate

## Usage

```text
ionosctl dataplatform nodepool create [flags]
```

## Aliases

For `dataplatform` command:

```text
[dp]
```

For `nodepool` command:

```text
[np]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Node Pool into an existing Data Platform Cluster. The Data Platform Cluster must be in state "ACTIVE" before creating a Node Pool. 

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* Cluster Id
* Name
* Node Count


## Options

```text
  -A, --annotations stringToString   Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host url (default "https://api.ionos.com")
  -z, --availability-zone string     The availability zone of the virtual datacenter region where the node pool resources should be provisioned. (default "AUTO")
      --cluster-id string            The unique ID of the Cluster (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [NodePoolId,Name,Version,NodeCount,DatacenterId,State])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                    The number of CPU cores per node.
      --cpu-family AUTO              A valid CPU family name or AUTO if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource. (default "AUTO")
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
  -L, --labels stringToString        Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
  -d, --maintenance-day string       Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -T, --maintenance-time string      Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -n, --name string                  The name for the  NodePool (required)
      --node-count int               The number of nodes that make up the node pool. (required)
  -o, --output string                Desired output format [text|json] (default "text")
  -q, --quiet                        Quiet output
      --ram string                   The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.
      --storage-size string          The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit
      --storage-type string          Storage Type
  -t, --timeout int                  Timeout option for waiting for NodePool to be in ACTIVE state[seconds] (default 600)
  -v, --verbose                      Print step-by-step process when running command
  -W, --wait-for-state               Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl dataplatform nodepool create --datacenter-id DATACENTER_ID --cluster-id CLUSTER_ID --name NAME --node-count NODE_COUNT
```

