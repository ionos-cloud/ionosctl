---
description: "Update a PostgreSQL Cluster"
---

# DbaasPostgresClusterUpdate

## Usage

```text
ionosctl dbaas postgres cluster update [flags]
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
  -u, --api-url string            Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cidr string               The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. e.g.: 192.168.1.100/24
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location State PostgresVersion Instances RAM Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance
      --datacenter-id string      The unique ID of the Datacenter to connect to your cluster. It has to be in the same location as the current datacenter
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int             The number of instances in your cluster. Minimum: 0. Maximum: 5
  -L, --lan-id string             The unique ID of the LAN to connect your cluster to
      --limit int                 Pagination limit: Maximum number of items to return per request (default 50)
      --maintenance-day string    Day of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59
  -n, --name string               The friendly name of your cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance. Size must be specified in multiples of 1024. The default unit is MB. Minimum: 4GB. e.g. --ram 4096, --ram 4096MB, --ram 4GB
      --remove-connection         Remove the connection completely
      --storage-size string       The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The PostgreSQL version of your cluster
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres cluster update -i CLUSTER_ID -n CLUSTER_NAME
```

