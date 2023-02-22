---
description: List the snapshots of your Mongo Cluster
---

# MongoSnapshotList

## Usage

```text
mongo snapshot list [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `snapshot` command:

```text
[snap backup snapshots backups]
```

For `list` command:

```text
[ls]
```

## Description

List the snapshots of your Mongo Cluster

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [SnapshotId CreationTime Size Version]
  -h, --help                help for list
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo cluster snapshots --cluster-id <cluster-id>
```

