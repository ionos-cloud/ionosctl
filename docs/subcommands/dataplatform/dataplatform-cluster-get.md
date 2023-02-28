---
description: Get a Dataplatform Cluster by ID
---

# DataplatformClusterGet

## Usage

```text
dataplatform cluster get [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `cluster` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Get a Dataplatform Cluster by ID

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -h, --help                help for get
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dataplatform cluster get --cluster-id <cluster-id>
```

