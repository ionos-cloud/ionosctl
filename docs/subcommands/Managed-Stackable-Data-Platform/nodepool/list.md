---
description: "List Dataplatform Nodepools of a certain cluster"
---

# DataplatformNodepoolList

## Usage

```text
ionosctl dataplatform nodepool list [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
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

List Dataplatform Nodepools of a certain cluster

## Options

```text
  -a, --all                 List all account nodepools, by iterating through all clusters first. May invoke a lot of GET calls
  -u, --api-url string      Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com/dataplatform")
  -i, --cluster-id string   The unique ID of the cluster. Must conform to the UUID format
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations ClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform nodepool list
```

