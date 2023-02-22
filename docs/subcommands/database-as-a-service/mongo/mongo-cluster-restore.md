---
description: Restore a Mongo Cluster by ID, using a snapshot
---

# MongoClusterRestore

## Usage

```text
mongo cluster restore [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `cluster` command:

```text
[c]
```

For `restore` command:

```text
[r]
```

## Description

Restore a Mongo Cluster by ID, using a snapshot

## Options

```text
  -i, --cluster-id string    The unique ID of the cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -h, --help                 help for restore
      --no-headers           When using text output, don't print headers
      --snapshot-id string   The unique ID of the snapshot you want to restore. (required)
```

## Examples

```text
ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>
```

