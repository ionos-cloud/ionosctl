---
description: Update a Data Platform NodePool
---

# DataplatformNodepoolUpdate

## Usage

```text
ionosctl dataplatform nodepool update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update the number of worker Nodes, to add labels, annotations, to update the maintenance day and time. 

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.
Required values to run command:

*  Cluster Id
*  NodePool Id

## Options

```text
  -A, --annotations stringToString   Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host url (default "https://api.ionos.com")
      --cluster-id string            The unique ID of the Cluster (required)
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [NodePoolId,Name,Version,NodeCount,DatacenterId,State])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
  -L, --labels stringToString        Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
  -d, --maintenance-day string       Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -T, --maintenance-time string      Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
      --node-count int               The number of nodes that make up the node pool.
  -i, --nodepool-id string           The unique ID of the Node Pool (required)
  -o, --output string                Desired output format [text|json] (default "text")
  -q, --quiet                        Quiet output
  -t, --timeout int                  Timeout option for waiting for NodePool to be in ACTIVE state[seconds] (default 600)
  -v, --verbose                      Print step-by-step process when running command
  -W, --wait-for-state               Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl dataplatform nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count NODE_COUNT
```

