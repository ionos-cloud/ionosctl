---
description: Delete a PostgreSQL Cluster
---

# DbaasPostgresClusterDelete

## Usage

```text
ionosctl dbaas postgres cluster delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified PostgreSQL Cluster from your account. You can wait for the cluster to be deleted with the wait-for-deletion option.

Required values to run command:

* Cluster Id

## Options

```text
  -a, --all                 Delete all Clusters
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -n, --name string         Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag
  -t, --timeout int         Timeout option for Cluster to be completely removed[seconds] (default 1200)
  -W, --wait-for-deletion   Wait for Cluster to be completely removed
```

## Examples

```text
ionosctl dbaas postgres cluster delete -i CLUSTER_ID
```

