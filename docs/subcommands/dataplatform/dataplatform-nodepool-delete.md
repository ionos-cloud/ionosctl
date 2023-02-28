---
description: Delete a Dataplatform Cluster by ID
---

# DataplatformNodepoolDelete

## Usage

```text
dataplatform nodepool delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete a Dataplatform Cluster by ID

## Options

```text
  -a, --all                  Delete all clusters. If cluster ID is provided, delete all nodepools in given cluster
      --cluster-id string    The unique ID of the cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations]
  -f, --force                Skip yes/no verification
  -h, --help                 help for delete
      --no-headers           When using text output, don't print headers
  -i, --nodepool-id string   The unique ID of the nodepool (required)
```

## Examples

```text
ionosctl dataplatform cluster delete --cluster-id <cluster-id>
```

