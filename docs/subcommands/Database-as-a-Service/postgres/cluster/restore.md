---
description: "Restore a PostgreSQL cluster in place from a backup"
---

# DbaasPostgresClusterRestore

## Usage

```text
ionosctl dbaas postgres cluster restore [flags]
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

For `restore` command:

```text
[r]
```

## Description

Trigger an in-place restore of an existing PostgreSQL cluster from one of its backups. This overwrites the cluster's current data with the state captured in the backup, so it is a destructive operation on the live cluster.

A backup is not a single frozen moment but a continuous recovery WINDOW. By default the backup is replayed in full (the latest available point). Pass --recovery-time to roll back to an earlier point in time inside that window (point-in-time recovery); the timestamp must fall within the backup's recovery window (see 'dbaas postgres backup get', field EarliestRecoveryTargetTime).

To instead create a NEW cluster from a backup while leaving the current one intact, use 'cluster create --backup-id'.

Required values to run command:

* Cluster Id
* Backup Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --backup-id string       ID of the backup to restore from. Completion lists only backups belonging to the chosen cluster (required)
  -i, --cluster-id string      The unique ID of the Cluster (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId DisplayName Location DatacenterId LanId Cidr Instances State PostgresVersion RAM Cores StorageSize StorageType MaintenanceWindow SynchronizationMode BackupLocation]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --recovery-time string   Replay the backup up to this point in time (ISO 8601 / RFC3339, e.g. 2024-01-15T10:00:00Z) for point-in-time recovery. Must fall within the backup's recovery window. If empty, the backup is applied in full
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Restore in place to the end of a backup
ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID

# Restore in place to a specific point in time within the backup's recovery window (PITR)
ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID --recovery-time 2024-01-15T10:00:00Z
```

