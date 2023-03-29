---
description: Get a PostgreSQL Cluster
---

# DbaasPostgresClusterGet

## Usage

```text
ionosctl dbaas postgres cluster get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a PostgreSQL Cluster by using its ID.

Required values to run command:

* Cluster Id

## Options

```text
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
      --no-headers          When using text output, don't print headers
  -t, --timeout int         Timeout option for Cluster to be in AVAILABLE state [seconds] (default 1200)
  -W, --wait-for-state      Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres cluster get -i CLUSTER_ID
```

