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

Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the Datacenter.

Required values to run command:

* Datacenter Id
* Lan Id
* IP
* Credentials for the User: Username and Password

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
  -b, --backup-id string          The unique ID of the backup you want to restore
      --cidr string               The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Example: 192.168.1.100/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                 The number of CPU cores per instance. Minimum: 1 (default 2)
  -D, --datacenter-id string      The unique ID of the Datacenter to connect to your cluster (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -I, --instances int             The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5 (default 1)
      --lan-id string             The unique ID of the LAN to connect your cluster to (required)
      --location-id string        The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation. If not set, it will be used Datacenter's location
  -d, --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
  -T, --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59
  -n, --name string               The friendly name of your cluster (default "UnnamedCluster")
  -o, --output string             Desired output format [text|json] (default "text")
  -P, --password string           Password for the initial postgres user (required)
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance. Size must be specified in multiples of 256. Minimum: 2048. The default unit is MB. e.g. --ram 2048 or --ram 2048MB (default "3GB")
  -R, --recovery-time string      If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely
      --storage-size string       The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB (default "20GB")
      --storage-type string       The storage type used in your cluster (default "HDD")
  -S, --sync string               Represents different modes of replication (default "ASYNCHRONOUS")
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -U, --username string           Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby) (required) (default "db-admin")
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The PostgreSQL version of your Cluster (default "13")
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl pg cluster create -V POSTGRES_VERSION --datacenter-id DATACENTER_ID --lan-id LAN_ID --ip IP_ADDRESS
```

