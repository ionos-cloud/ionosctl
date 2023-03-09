---
description: Get a Mongo Cluster by ID
---

# MongoClusterGet

## Usage

```text
mongo cluster get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get a Mongo Cluster by ID

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -h, --help                help for get
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo cluster get --cluster-id <cluster-id>
```

