---
description: Get Dataplatform Nodepool by cluster and nodepool id
---

# DataplatformNodepoolGet

## Usage

```text
dataplatform nodepool get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get Dataplatform Nodepool by cluster and nodepool id

## Options

```text
      --cluster-id string    The unique ID of the cluster. Must conform to the UUID format
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations]
  -h, --help                 help for get
      --no-headers           When using text output, don't print headers
  -i, --nodepool-id string   The unique ID of the nodepool. Must conform to the UUID format
```

## Examples

```text
ionosctl dataplatform nodepool get
```

