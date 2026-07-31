---
description: "Update a PostgreSQL cluster"
---

# DbaasPostgresClusterUpdate

## Usage

```text
ionosctl dbaas postgres cluster update [flags]
```

## Aliases

For `postgres` command:

```text
[pg pgsql postgresql psql]
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

Update attributes of an existing PostgreSQL cluster. Only the flags you pass are changed; everything else is left untouched (a PATCH).

You can scale compute (--cores, --ram) and grow storage (--storage-size; storage can only be increased, not shrunk), rename the cluster (--name), upgrade the engine (--version), adjust the maintenance window (--maintenance-day and --maintenance-time must be given together), or change the network connection (--datacenter-id, --lan-id, --cidr).

--storage-type and --location cannot be changed after creation. Use --remove-connection to detach the cluster from its LAN entirely.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cidr string               New IP address and subnet for the cluster on the LAN, in CIDR notation (e.g. 192.168.1.100/24). Must not overlap the reserved ranges 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location DatacenterId LanId Cidr Instances State PostgresVersion RAM Cores StorageSize StorageType MaintenanceWindow SynchronizationMode BackupLocation]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 New number of CPU cores per instance. Leave at 0 to keep the current value
      --datacenter-id string      Move the cluster's connection to this virtual datacenter. It must be in the same location as the current datacenter
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int             New total number of instances (1 master + n-1 standbys). Maximum: 5. Leave at 0 to keep the current count
  -L, --lan-id string             Move the cluster's connection to this LAN (within the target datacenter)
      --limit int                 Maximum number of items to return per request (default 50)
      --maintenance-day string    New day of the week (e.g. Monday) for the weekly 4-hour maintenance window. Must be set together with --maintenance-time
      --maintenance-time string   New start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window. Must be set together with --maintenance-day
  -n, --name string               New human-friendly display name for the cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                New memory per instance. Must be a multiple of 1024MB and at least 4GB. Default unit is MB. e.g. --ram 4096, --ram 4096MB, --ram 4GB
      --remove-connection         Detach the cluster from its LAN, removing the network connection entirely. Mutually exclusive with setting --datacenter-id/--lan-id/--cidr
      --storage-size string       New storage per instance. Storage can only be increased, never decreased. Default unit is MB. e.g.: --storage-size 20480, --storage-size 20480MB, --storage-size 20GB
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            Upgrade the PostgreSQL engine to this major version (e.g. 16). Only forward upgrades are supported; the cluster is briefly unavailable during the upgrade
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a cluster
ionosctl dbaas postgres cluster update -i CLUSTER_ID --name new-name

# Scale up compute and storage, and set a maintenance window (day and time must be given together)
ionosctl dbaas postgres cluster update -i CLUSTER_ID --cores 4 --ram 16GB --storage-size 200GB --maintenance-day Saturday --maintenance-time 04:00:00
```

