---
description: "Create a MariaDB Cluster"
---

# DbaasMariadbV2ClusterCreate

## Usage

```text
ionosctl dbaas mariadb-v2 cluster create [flags]
```

## Aliases

For `mariadb-v2` command:

```text
[maria-v2 mar-v2 ma-v2]
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

Use this command to create a new MariaDB Cluster. The mode is determined by the number of instances: one instance is standalone, everything else is a replication with one primary and n-1 secondary instances.

There are two ways to create a cluster, both requiring the same connection and credential flags (--datacenter-id, --lan-id, --cidr, --user, --password, --database) plus --location:
  1. Empty cluster: pass --version (defaults to a supported version) and sizing flags (--instances, --cores, --ram, --storage-size).
  2. From a backup: additionally pass --backup-id. The cluster version is taken from the backup (so --version is not needed; if given, it must match the backup's version). Optionally pass --recovery-time to restore to a point in time within the backup's window.

## Options

```text
  -u, --api-url string                Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
      --backup-id string              Create the cluster from this backup instead of empty. The connection/credential flags and --location are still required; the cluster version is taken from the backup
  -B, --backup-location string        The Object Storage location where backups will be stored. For added data safety, use a different location than the cluster (default "eu-central-4")
      --backup-retention-days int32   Configures how many days cluster backups are retained before being automatically deleted. Minimum: 1, Maximum: 365 (default 30)
      --cidr string                   The IP and subnet for the cluster. All IPs must be in a /24 network. e.g.: 192.168.1.100/24 (required)
      --cols strings                  Set of columns to be printed on output 
                                      Available columns: [ClusterId Name DnsName Version Instances State Cores Ram StorageSize Description BackupLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr StatusMessage]
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32                   The number of CPU cores per instance (default 1)
      --database string               The name of the initial database created for the user (required)
      --datacenter-id string          The unique ID of the Datacenter to connect to your cluster. Must be in the same location as the cluster (required)
  -D, --depth int                     Level of detail for response objects (default 1)
  -F, --filters strings               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
      --instances int32               The total number of instances in the cluster (one primary and n-1 secondary). For a standalone instance, use 1 (default 1)
      --lan-id string                 The unique ID of the LAN to connect your cluster to (required)
      --limit int                     Maximum number of items to return per request (default 50)
  -l, --location string               Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci
      --logs-enabled                  Enable collection and reporting of logs for this cluster
      --maintenance-day string        Day of the week for the MaintenanceWindow. Defaults to a random day during Mon-Fri. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string       Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --metrics-enabled               Enable collection and reporting of metrics for this cluster
  -n, --name string                   The friendly name of your cluster (default "UnnamedCluster")
      --no-headers                    Don't print table headers when table output is used
      --offset int                    Number of items to skip before starting to collect the results
      --order-by string               Property to order the results by
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
      --password string               Password for the initial user (required)
      --query string                  JMESPath query string to filter the output
  -q, --quiet                         Quiet output
      --ram string                    The amount of memory per instance. e.g. --ram 4, --ram 4GB. Minimum 4GB (required) (default "4GB")
  -R, --recovery-time string          Advanced: with --backup-id, restore to a specific point in time WITHIN the backup's recovery window (PITR), as an ISO 8601 timestamp. Defaults to the latest point in the window
      --storage-size string           The size of the storage per instance. e.g. --storage-size 10 or --storage-size 10GB (default "10GB")
  -t, --timeout int                   Timeout in seconds for --wait and other wait operations (default 600)
      --user string                   The initial username (required)
  -v, --verbose count                 Increase verbosity level [-v, -vv, -vvv]
      --version string                The MariaDB version of your cluster. Ignored when --backup-id is set (the backup's version is used) (required) (default "10.11")
  -w, --wait                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an empty cluster
ionosctl dbaas mariadb-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --version <version>

# Create a cluster from an existing backup (version is taken from the backup)
ionosctl dbaas mariadb-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --backup-id <backup-id>
```

