---
description: "Restore a Mongo Cluster in place from one of its snapshots"
---

# DbaasMongoSnapshotRestore

## Usage

```text
ionosctl dbaas mongo snapshot restore [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `snapshot` command:

```text
[snap backup snapshots backups]
```

For `restore` command:

```text
[r]
```

## Description

Roll a MongoDB cluster back in place to the state captured by one of its own snapshots. This overwrites the current data of the cluster identified by --cluster-id with the contents of --snapshot-id (list a cluster's snapshots with `snapshot list --cluster-id <id>`).

How snapshots accumulate: an initial snapshot is taken when the cluster is created (the initial sync, usually within 24h); another is created after each restore; thereafter a base snapshot is taken every 24h and a full snapshot every Sunday. Snapshots are retained for the last 7 days, so recovery is possible up to a week back. Playground clusters have no snapshots and cannot be restored.

Constraints:
  - You can only restore from a snapshot whose MongoDB version is the same as, or older (by patch) than, the cluster's current version.
  - Snapshots are stored in IONOS S3 Object Storage in the cluster's region (eu-central-2 where S3 is unavailable).

Enterprise clusters additionally support point-in-time recovery WITHIN a snapshot's window via the API's recoveryTargetTime; this CLI restore targets a whole snapshot.

## Options

```text
  -u, --api-url string                                Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string                             The unique ID of the cluster to restore in place (its current data will be overwritten) (required)
      --cols strings                                  Set of columns to be printed on output 
                                                      Available columns: [SnapshotId CreationTime Size Version]
  -c, --config string                                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                                     Level of detail for response objects (default 1)
  -F, --filters strings                               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                         Force command to execute without user input
  -h, --help                                          Print usage
      --limit int                                     Maximum number of items to return per request (default 50)
      --no-headers                                    Don't print table headers when table output is used
      --offset int                                    Number of items to skip before starting to collect the results
      --order-by string                               Property to order the results by
  -o, --output string                                 Desired output format [text|json|api-json] (default "text")
      --query string                                  JMESPath query string to filter the output
  -q, --quiet                                         Quiet output
      --snapshot-id snapshot list --cluster-id <id>   The snapshot to restore from. Must belong to this cluster and its MongoDB version must be the same or older than the cluster's. List options with snapshot list --cluster-id <id> (required)
  -t, --timeout int                                   Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                                 Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>
```

