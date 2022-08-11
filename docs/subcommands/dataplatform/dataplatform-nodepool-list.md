---
description: List Data Platform NodePools
---

# DataplatformNodepoolList

## Usage

```text
ionosctl dataplatform nodepool list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of all contained NodePools in a selected Data Platform Cluster.

Required values to run command:

*  Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [NodePoolId,Name,Version,NodeCount,DatacenterId,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform nodepool list --cluster-id CLUSTER_ID
```

