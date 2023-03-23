---
description: List Dataplatform Nodepools of a certain cluster
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster. Must conform to the UUID format
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations]
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
ionosctl dataplatform nodepool list
```

