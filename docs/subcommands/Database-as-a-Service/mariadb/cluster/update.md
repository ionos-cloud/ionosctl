---
description: "Update a MariaDB Cluster"
---

# DbaasMariadbClusterUpdate

## Usage

```text
ionosctl dbaas mariadb cluster update [flags]
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

For `update` command:

```text
[u]
```

## Description

Update mutable attributes of an existing MariaDB cluster. Only the flags you pass are changed; everything else is left untouched. The cluster must be in state AVAILABLE for the update to be accepted.

Some changes are one-directional and cannot be reverted:
  - --version can only be upgraded (10.6 -> 10.11), never downgraded.
  - --instances can only be increased, and must stay odd (1 -> 3 -> 5).
  - --storage-size can only be increased, never shrunk.
  - --cores and --ram can be scaled up or down.

--maintenance-day and --maintenance-time must be supplied together (a maintenance window has both a day and a start time). The connection (datacenter/LAN/CIDR) and initial credentials are set at creation and cannot be changed here.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
  -i, --cluster-id string         The unique ID of the cluster to update. Must be in state AVAILABLE (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               New CPU core count per instance (minimum 1). Can be scaled up or down
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           New instance count (primary + secondaries). Can only be increased and must stay odd: 1, 3 or 5. Adding instances converts a standalone cluster into a high-availability replica set
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --maintenance-day string    New day of the week for the weekly maintenance window, e.g. Monday. Must be supplied together with --maintenance-time
      --maintenance-time string   New start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window. Must be supplied together with --maintenance-day
  -n, --name string               New human-friendly display name for the cluster (max 63 characters)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                New memory per instance, e.g. --ram 8GB. Minimum 4GB, whole GB only. Can be scaled up or down
      --storage-size string       New storage per instance, e.g. --storage-size 200GB. Can only be increased (never shrunk), up to 2000GB (2TB)
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            Upgrade the MariaDB version (one of: 10.6, 10.11). Upgrade only; downgrades are rejected by the API
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Upgrade the MariaDB version
ionosctl dbaas mariadb cluster update --cluster-id <cluster-id> --version 10.11

# Scale compute and move the maintenance window (day and time must be given together)
ionosctl dbaas mariadb cluster update --cluster-id <cluster-id> --cores 8 --ram 32GB --storage-size 200GB --maintenance-day Saturday --maintenance-time 02:00:00
```

