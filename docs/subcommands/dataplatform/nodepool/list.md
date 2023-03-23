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
  -i, --cluster-id string   The unique ID of the cluster. Must conform to the UUID format
```

## Examples

```text
ionosctl dataplatform nodepool list
```

