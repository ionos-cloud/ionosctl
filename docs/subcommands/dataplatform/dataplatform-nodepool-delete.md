---
description: Delete a Data Platform NodePool
---

# DataplatformNodepoolDelete

## Usage

```text
ionosctl dataplatform nodepool delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

This command deletes a Data Platform Node Pool within an existing Data Platform Cluster.

Required values to run command:

*  Cluster Id
*  NodePool Id

## Options

```text
  -a, --all                  Delete all the Data Platform Node Pools within an existing Data Platform Cluster.
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique ID of the Cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [NodePoolId,Name,Version,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --nodepool-id string   The unique ID of the Node Pool (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

