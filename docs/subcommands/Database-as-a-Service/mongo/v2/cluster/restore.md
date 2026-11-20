---
description: "Restore a MongoDB Cluster in place to a point in time"
---

# DbaasMongoV2ClusterRestore

## Usage

```text
ionosctl dbaas mongo-v2 cluster restore [flags]
```

## Aliases

For `mongo-v2` command:

```text
[mongodb-v2 mg-v2]
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

Use this command to trigger an in-place restore of the specified MongoDB Cluster from its own snapshots. The cluster's current data is overwritten with the restored data, and the cluster may be briefly unavailable while this happens.

Snapshots are not single points in time — each keeps the cluster's oplog, forming a continuous recovery WINDOW (see a snapshot's earliestRecoveryTargetTime / latestRecoveryTargetTime via `snapshot get`). `--recovery-time` picks the exact moment inside that window to restore to. Omit `--recovery-time` (or use `now`) to restore to the latest point. Accepted formats: `now`, a date (`2025-01-02`), a date-time (`"2025-01-02 15:00"`), or a full RFC3339 timestamp (`2025-01-02T15:00:00Z`); values without a timezone are treated as UTC. Tab-completion suggests the window's boundaries.

To instead create a NEW cluster from a specific snapshot, use `cluster create --snapshot-id`.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mongodbv2' and env var 'IONOS_API_URL' (default "https://mongodb.%s.ionos.com")
  -i, --cluster-id string      The unique ID of the Cluster (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId Name Edition Version Instances State ConnectionUri Cores Ram StorageSize TemplateId ShardingEnabled Shards BiConnectorEnabled Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username Database StatusMessage]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. When unset, list commands query all locations. Can be one of: gb/bhx, us/ewr. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -R, --recovery-time string   Point inside the recovery window to restore to: 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). Defaults to the latest point (default "now")
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo-v2 cluster restore --location <location> --cluster-id <cluster-id> --recovery-time 2025-01-02T15:00:00Z
```

