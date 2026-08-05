---
description: "Restore a MariaDB Cluster in place to a point in time"
---

# DbaasMariadbV2ClusterRestore

## Usage

```text
ionosctl dbaas mariadb-v2 cluster restore [flags]
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

For `restore` command:

```text
[r]
```

## Description

Use this command to trigger an in-place restore of the specified MariaDB Cluster from its own backups.

Backups are not single points in time — they form a continuous recovery WINDOW (see a backup's earliestRecoveryTargetTime via `backup get`). `--recovery-time` zooms into a specific moment inside that window; the cluster is rolled back to the nearest point at or before it. Omit `--recovery-time` (or use `now`) to restore to the latest point. Accepted formats: `now`, a date (`2025-01-02`), a date-time (`"2025-01-02 15:00"`), or a full RFC3339 timestamp (`2025-01-02T15:00:00Z`); values without a timezone are treated as UTC.

To instead create a NEW cluster from a specific backup, use `cluster create --backup-id`.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
  -i, --cluster-id string      The unique ID of the Cluster (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId Name DnsName Version Instances State Cores Ram StorageSize Description BackupLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr StatusMessage]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --database string        Database for the credentials. Defaults to the cluster's current database
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --password string        Password for the database user. Required because the API does not return it on GET requests (minimum length 10) (required)
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -R, --recovery-time string   Point inside the recovery window to restore to: 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). The nearest point at or before this time is used; defaults to the latest (default "now")
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
      --user string            Username for the database user. Defaults to the cluster's current username
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mariadb-v2 cluster restore --cluster-id <cluster-id> --password <password> --recovery-time 2025-01-02T15:00:00Z
```

