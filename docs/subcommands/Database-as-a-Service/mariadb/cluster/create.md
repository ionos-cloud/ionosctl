---
description: "Create a DBaaS MariaDB cluster"
---

# DbaasMariadbClusterCreate

## Usage

```text
ionosctl dbaas mariadb cluster create [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar ma]
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

Create a new managed MariaDB cluster.

A cluster is a set of instances running the same MariaDB version: one primary (read-write) and, for a replicated setup, n-1 secondaries (read-only replicas kept in sync with the primary). The instance count must be odd so a quorum can elect a primary: use --instances 1 for a single standalone instance (no high availability), or 3 or 5 for a high-availability replica set. Every instance is sized identically via --cores, --ram and --storage-size.

The cluster is reached only over a private LAN inside one of your VDCs, never over the public internet. --datacenter-id, --lan-id and --cidr together define that connection; the datacenter must be in the same location (region) as the cluster, exactly one connection is allowed, and the --cidr address must sit inside the chosen LAN's subnet. After creation the cluster is reachable at the DNS name shown in 'cluster get'.

--user and --password create the initial database user; the password is write-only and is never returned by the API afterwards, so store it safely. A newly created cluster starts in state BUSY and must reach AVAILABLE before you can connect to it or run further operations (update, restore) against it.

To provision a cluster from an existing backup instead of empty, use the API's fromBackup option; in ionosctl, inspect available backups with 'ionosctl dbaas mariadb backup list'.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
      --cidr string               The cluster's IP and subnet within the LAN in CIDR notation, e.g. 192.168.1.100/24 (use a /24 network). The address must lie within the LAN's subnet. Unavailable ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               Number of CPU cores per instance (minimum 1). Applies to every instance in the cluster (default 1)
      --datacenter-id string      ID of the Virtual Data Center (VDC) hosting the private LAN the cluster connects to. Must be in the same location (region) as the cluster (required)
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           Number of instances in the cluster: 1 (standalone, no high availability) or an odd number (3 or 5) for a primary + secondaries replica set. Must be odd so the replicas can elect a primary. Range 1-5 (default 1)
      --lan-id string             Numeric ID of the private LAN (inside --datacenter-id) the cluster attaches to. The cluster is reachable only over this LAN, not the public internet (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --maintenance-day string    Day of the week the weekly maintenance window opens (e.g. Monday). Pairs with --maintenance-time. Defaults to a random weekday (Mon-Fri) (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window during which IONOS may apply patches and minor upgrades. Pairs with --maintenance-day. Defaults to a random time between 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               Human-friendly display name for the cluster (max 63 characters). Not a DNS name and need not be unique (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           Password for the initial database user (10-63 characters). Write-only: the API never returns it afterwards, so record it securely (required)
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                Memory per instance, e.g. --ram 4GB. Minimum 4GB; must be a whole number of GB. The upper bound is set by your contract quota (default "4GB")
      --storage-size string       Storage per instance, e.g. --storage-size 10 or --storage-size 10GB. Minimum 10GB, maximum 2000GB (2TB). Can later be increased (never shrunk) via 'cluster update' (default "10")
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               Username for the initial database user (1-16 chars, must start with a letter, letters/digits/underscores only). Reserved names such as mariadb, admin and standby are rejected (required)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The MariaDB server version to run. One of: 10.6, 10.11. Can later be upgraded (never downgraded) via 'cluster update' (required) (default "10.6")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a standalone MariaDB cluster (MariaDB 10.6, 1 instance)
i db mariadb cluster create --name NAME --version VERSION --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --user USER --password PASSWORD 

# Create a high-availability cluster (3 instances) with explicit sizing and a fixed maintenance window
ionosctl dbaas mariadb cluster create --name prod-db --version 10.11 --instances 3 --cores 4 --ram 16GB --storage-size 100GB --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr 192.168.1.100/24 --user cluster_admin --password <password> --maintenance-day Sunday --maintenance-time 03:00:00
```

