---
description: Update a PostgreSQL Cluster
---

# PgClusterUpdate

## Usage

```text
ionosctl pg cluster update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update attributes of a PostgreSQL Cluster.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location BackupEnabled LifecycleStatus PostgresVersion Replicas RamSize CpuCoreCount StorageSize StorageType DatacenterId LanId IpAddress MaintenanceWindow] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,IpAddress,Replicas,LifecycleStatus])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-core-count int        The number of CPU cores per replica
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -d, --maintenance-day string    WeekDay for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
  -T, --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59
  -n, --name string               The friendly name of your cluster
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --ram-size string           The amount of memory per replica in IEC format. Value must be a multiple of 1024Mi and at least 2048Mi
  -R, --replicas int              The number of replicas in your cluster
      --storage-size string       The amount of storage per replica. It is expected IEC format like 2Gi or 500Mi
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The PostgreSQL version of your cluster
```

## Examples

```text
ionosctl pg cluster update -i CLUSTER_ID -n CLUSTER_NAME
```

