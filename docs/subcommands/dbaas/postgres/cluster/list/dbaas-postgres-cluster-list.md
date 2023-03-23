---
description: List PostgreSQL Clusters
---

# DbaasPostgresClusterList

## Usage

```text
ionosctl dbaas postgres cluster list [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
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

Use this command to retrieve a list of PostgreSQL Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
      --cols strings   Set of columns to be printed on output 
                       Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -n, --name string    Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers     When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas postgres cluster list
```

