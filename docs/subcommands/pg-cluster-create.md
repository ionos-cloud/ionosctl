---
description: Create a PostgreSQL Cluster
---

# PgClusterCreate

## Usage

```text
ionosctl pg cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the VDC (VirtualDataCenter), the unique ID of the LAN. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the VDC.

Required values to run command:

* Datacenter Id
* Lan Id
* IP

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
  -b, --backup-id string          The unique ID of the backup you want to restore
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location BackupEnabled LifecycleStatus PostgresVersion Replicas RamSize CpuCoreCount StorageSize StorageType DatacenterId LanId IpAddress MaintenanceWindow] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,IpAddress,Replicas,LifecycleStatus])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-core-count int        The number of CPU cores per replica (default 4)
  -D, --datacenter-id string      The unique ID of the VDC to connect to your cluster (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --ip string                 The private IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Example: 192.168.1.100/24 (required)
      --lan-id string             The unique Lan ID (required)
      --location-id string        The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests). If not set, it will be used VDC's location (default "de/fra")
  -d, --maintenance-day string    WeekDay for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
  -T, --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59
  -n, --name string               The friendly name of your cluster (default "UnnamedCluster")
  -o, --output string             Desired output format [text|json] (default "text")
      --password string           Password for the database user to be created (default "password")
  -q, --quiet                     Quiet output
      --ram-size string           The amount of memory per replica in IEC format. Value must be a multiple of 1024Mi and at least 2048Mi (default "2Gi")
  -R, --replicas int              The number of replicas in your cluster. Minimum: 1. Maximum: 5 (default 1)
      --storage-size string       The amount of storage per replica. It is expected IEC format like 2Gi or 500Mi (default "20Gi")
      --storage-type string       The storage type used in your cluster (default "HDD")
  -S, --sync string               Represents different modes of replication (default "asynchronous")
      --time string               If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
      --username string           Username for the database user to be created. Some system usernames are restricted (e.g. postgres, admin, standby) (default "db-admin")
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The PostgreSQL version of your Cluster (default "13")
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl pg cluster create -V POSTGRES_VERSION --datacenter-id DATACENTER_ID --lan-id LAN_ID --ip IP_ADDRESS
```

