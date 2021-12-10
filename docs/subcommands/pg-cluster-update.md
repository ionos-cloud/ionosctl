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
  -C, --cidr string               The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. e.g.: 192.168.1.100/24
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                 The number of CPU cores per instance
  -D, --datacenter-id string      The unique ID of the Datacenter to connect to your cluster. It has to be in the same location as the current datacenter
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -I, --instances int             The number of instances in your cluster. Minimum: 0. Maximum: 5
  -L, --lan-id string             The unique ID of the LAN to connect your cluster to
  -d, --maintenance-day string    Day of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
  -T, --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59
  -n, --name string               The friendly name of your cluster
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance. Size must be specified in multiples of 1024. The default unit is MB. Minimum: 2048. e.g. --ram 2048, --ram 2048MB, --ram 2GB
      --remove-connection         Remove the connection completely
      --storage-size string       The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The PostgreSQL version of your cluster
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl pg cluster update -i CLUSTER_ID -n CLUSTER_NAME
```

