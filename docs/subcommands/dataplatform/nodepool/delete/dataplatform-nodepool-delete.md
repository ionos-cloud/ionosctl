---
description: Delete a Dataplatform Cluster by ID
---

# DataplatformNodepoolDelete

## Usage

```text
ionosctl dataplatform nodepool delete [flags]
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
  -f, --force                Skip yes/no verification
  -i, --nodepool-id string   The unique ID of the nodepool (required)
```

## Examples

```text
ionosctl dataplatform cluster delete --cluster-id <cluster-id>
```

