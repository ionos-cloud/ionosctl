---
description: List Dataplatform Nodepools of a certain cluster
---

# DataplatformNodepoolList

## Usage

```text
dataplatform nodepool list [flags]
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
  -i, --cluster-id string   The unique ID of the cluster. Must conform to the UUID format
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations]
  -h, --help                help for list
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dataplatform nodepool list
```

