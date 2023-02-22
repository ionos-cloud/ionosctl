---
description: List Mongo Clusters
---

# MongoClusterList

## Usage

```text
mongo cluster list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
      --cols strings   Set of columns to be printed on output 
                       Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -h, --help           help for list
  -n, --name string    Response filter to list only the Mongo Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers     When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo cluster list
```

