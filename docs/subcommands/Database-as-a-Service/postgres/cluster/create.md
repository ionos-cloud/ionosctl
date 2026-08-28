---
description: "Create a PostgreSQL cluster"
---

# DbaasPostgresClusterCreate

## Usage

```text
ionosctl dbaas postgres cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create a new PostgreSQL cluster.

The cluster attaches to a private LAN inside one of your virtual datacenters (VDC), so at minimum you must supply the network connection (--datacenter-id, --lan-id, --cidr) and the credentials for the initial database user (--db-username, --db-password). Every other property has a sensible default (PostgreSQL 15, 1 instance, 2 cores, 4GB RAM, 20GB HDD storage, ASYNCHRONOUS replication).

--location is the physical region the instances live in and cannot be changed later; if omitted it is inherited from the datacenter's location.

Sizing constraints: --instances 1-5 (1 master + n-1 read-standbys), --cores min 1, --ram min 4GB (must be a multiple of 1024MB), --storage-size 10GB-2TB.

To seed the cluster from an existing backup instead of an empty database (a clone), add --backup-id, and optionally --recovery-time to replay to a point in time within that backup's recovery window.

Provisioning is asynchronous: the cluster returns immediately in a non-AVAILABLE state. Wait for AVAILABLE (e.g. via 'cluster get') before creating additional databases or users.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* Credentials for the database user: Username and Password

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --backup-id string          Seed the new cluster from this backup instead of an empty database (clone). The backup's PostgreSQL version must be compatible with --version. See 'dbaas postgres backup list'
      --backup-location string    Object Storage (S3) region where automated backups are stored: de, eu-south-2, eu-central-3, eu-central-4, us-central-1. Defaults to a region derived from the cluster location. Cannot be changed after creation
  -C, --cidr string               IP address and subnet the master reserves on the LAN, in CIDR notation (e.g. 192.168.1.100/24). Must not overlap the reserved ranges 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName Location DatacenterId LanId Cidr Instances State PostgresVersion RAM Cores StorageSize StorageType MaintenanceWindow SynchronizationMode BackupLocation]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 Number of CPU cores allocated to each instance. Minimum: 1 (default 2)
      --datacenter-id string      ID of the virtual datacenter (VDC) hosting the LAN the cluster attaches to. Its region also sets the default --location (required)
      --db-password string        Password for the initial user. Minimum 10 characters (required)
      --db-username string        Username of the initial PostgreSQL superuser-equivalent role created with the cluster. Reserved system names such as postgres, admin and standby are not allowed (required)
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int             Total number of PostgreSQL instances: 1 master plus n-1 read-standby replicas. 1 means a single standalone instance (no HA). Minimum: 1, Maximum: 5 (default 1)
  -L, --lan-id string             Numeric ID of the LAN, within the chosen datacenter, that the cluster connects to. The cluster is reachable only from this private LAN (required)
      --limit int                 Maximum number of items to return per request (default 50)
      --location string           Physical region where the instances are provisioned (e.g. de/fra, de/txl, gb/lhr, us/las). Cannot be modified after creation. If omitted, the datacenter's location is used, so it must match the datacenter's region
      --maintenance-day string    Day of the week (e.g. Monday) for the weekly 4-hour maintenance window. Set together with --maintenance-time
      --maintenance-time string   Start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window during which the service may apply updates. Set together with --maintenance-day. If omitted, a window is assigned automatically
  -n, --name string               Human-friendly display name for the cluster (does not have to be unique) (default "UnnamedCluster")
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                Memory per instance. Must be a multiple of 1024MB and at least 4GB. Default unit is MB if none is given. e.g. --ram 4096, --ram 4096MB, --ram 4GB (default "4GB")
      --recovery-time string      Only with --backup-id: replay the backup up to this point in time (ISO 8601 / RFC3339, e.g. 2024-01-15T10:00:00Z) for point-in-time recovery. Must fall within the backup's recovery window. If empty, the backup is applied in full (latest available point)
      --storage-size string       Disk storage per instance. Default unit is MB if none is given. Minimum 10GB, maximum 2TB. e.g.: --storage-size 20480, --storage-size 20480MB, --storage-size 20GB (default "20GB")
      --storage-type string       Disk performance tier: HDD (spinning disk, cheapest), SSD_STANDARD (general-purpose SSD), SSD_PREMIUM (highest IOPS). SSD is deprecated and treated as SSD_PREMIUM. Cannot be changed after creation (default "HDD")
  -S, --sync string               Replication mode between master and standbys. ASYNCHRONOUS: fastest, standbys may lag and a failover can lose the last transactions. STRICTLY_SYNCHRONOUS: a write is only acknowledged once a standby has it, safest but slower. SYNCHRONOUS is deprecated; prefer ASYNCHRONOUS or STRICTLY_SYNCHRONOUS (default "ASYNCHRONOUS")
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The major PostgreSQL engine version to run (e.g. 13, 14, 15, 16). See 'dbaas postgres version list' for the versions currently offered (default "15")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Basic: smallest cluster with defaults (PostgreSQL 15, 1 instance, 2 cores, 4GB RAM, 20GB HDD, ASYNCHRONOUS). Location is inherited from the datacenter
ionosctl dbaas postgres cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --db-username dbadmin --db-password 'S3cr3tPassw0rd'

# Advanced: named 3-instance HA cluster, sized, on SSD Premium, with a maintenance window and a synchronous replication mode
ionosctl dbaas postgres cluster create --name prod-orders --version 16 --instances 3 --cores 4 --ram 8GB --storage-size 100GB --storage-type SSD_PREMIUM --sync-mode STRICTLY_SYNCHRONOUS --location de/fra --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --backup-location eu-central-3 --maintenance-day Sunday --maintenance-time 03:00:00 --db-username dbadmin --db-password 'S3cr3tPassw0rd'

# Clone: create a new cluster from an existing backup, replayed to a point in time (PITR)
ionosctl dbaas postgres cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --db-username dbadmin --db-password 'S3cr3tPassw0rd' --backup-id BACKUP_ID --recovery-time 2024-01-15T10:00:00Z
```

