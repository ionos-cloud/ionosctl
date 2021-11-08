---
description: Delete a PostgreSQL Cluster
---

# PgClusterDelete

## Usage

```text
ionosctl pg cluster delete [flags]
```

## Aliases

For `pg` command:

```text
[postgres]
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

Use this command to delete a specified PostgreSQL Cluster from your account.

Required values to run command:

* Cluster Id

## Options

```text
  -a, --all                 Delete all Clusters
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName Location BackupEnabled LifecycleStatus PostgresVersion Replicas RamSize CpuCoreCount StorageSize StorageType DatacenterId LanId IpAddress MaintenanceWindow] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,IpAddress,Replicas,LifecycleStatus])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -n, --name string         Delete all Clusters after filtering based on name. Can be used with --all flag
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl pg cluster delete -i CLUSTER_ID
```

