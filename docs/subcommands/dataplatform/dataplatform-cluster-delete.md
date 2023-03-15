---
description: Delete a Dataplatform Cluster by ID
---

# DataplatformClusterDelete

## Usage

```text
dataplatform cluster delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete a Dataplatform Cluster by ID

## Options

```text
  -a, --all                 Delete all clusters
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -f, --force               Skip yes/no verification
  -h, --help                help for delete
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dataplatform cluster delete --cluster-id <cluster-id>
```

