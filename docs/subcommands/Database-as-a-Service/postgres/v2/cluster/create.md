---
description: "Create a PostgreSQL Cluster"
---

# DbaasPostgresV2ClusterCreate

## Usage

```text
ionosctl dbaas postgres-v2 cluster create [flags]
```

## Aliases

For `postgres-v2` command:

```text
[pg-v2]
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

Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* Credentials for the database user: Username and Password

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'todo' and env var 'IONOS_API_URL' (default "https://postgresql.%s.ionos.com")
  -b, --backup-id string          The unique ID of the backup you want to restore from when creating this cluster
  -B, --backup-location string    The S3 location where the backups will be stored. Defaults to 'de' (default "de")
  -C, --cidr string               The IP and subnet for the cluster. Note the following unavailable IP range: 10.208.0.0/12. e.g.: 192.168.1.100/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName DnsName PostgresVersion Instances Ram Cores StorageSize State SyncMode MaintenanceDay MaintenanceTime BackupLocation DatacenterId LanId Cidr] (default [ClusterId,DisplayName,DnsName,PostgresVersion,Instances,Ram,Cores,StorageSize,State,SyncMode])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance. Minimum: 1 (default 2)
      --datacenter-id string      The unique ID of the Datacenter to connect to your cluster (required)
  -P, --db-password string        Password for the initial postgres user (required)
  -U, --db-username string        Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby) (required)
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -I, --instances int             The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5 (default 1)
  -L, --lan-id string             The unique ID of the LAN to connect your cluster to (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, gb/bhx, us/las, us/mci, us/ewr (default "de/txl")
  -d, --maintenance-day string    Day of the week for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. Defaults to a random day during Mon-Fri (default "Random (Mon-Fri 10:00-16:00)")
  -T, --maintenance-time string   Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00 (default "14:00:00")
  -n, --name string               The friendly name of your cluster (default "UnnamedCluster")
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance in GB. e.g. --ram 4096, --ram 4096MB, --ram 4GB (default "4GB")
  -R, --recovery-time string      If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely
      --storage-size string       The amount of storage per instance in GB. e.g.: --storage-size 20480 or --storage-size 20480MB or --storage-size 20GB (default "20GB")
  -S, --sync-mode string          Replication mode. Represents different modes of replication: ASYNCHRONOUS, SYNCHRONOUS, STRICTLY_SYNCHRONOUS (default "ASYNCHRONOUS")
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -V, --version string            The PostgreSQL version of your Cluster (default "15")
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres-v2 cluster create --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --db-username <username> --db-password <password>
```

