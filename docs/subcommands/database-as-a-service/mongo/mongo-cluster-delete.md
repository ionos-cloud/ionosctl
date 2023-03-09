---
description: Delete a Mongo Cluster by ID
---

# MongoClusterDelete

## Usage

```text
mongo cluster delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete a Mongo Cluster by ID

## Options

```text
  -a, --all                 Delete all mongo clusters
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -f, --force               Skip yes/no verification
  -h, --help                help for delete
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo cluster delete --cluster-id <cluster-id>
```

