---
description: "Restore a Replica Set from a Snapshot"
---

# DbaasInMemoryDbSnapshotRestoreCreate

## Usage

```text
ionosctl dbaas in-memory-db snapshot restore create [flags]
```

## Aliases

For `snapshot` command:

```text
[snaps snap backup backups snapshots]
```

For `restore` command:

```text
[restores backup backups]
```

For `create` command:

```text
[c post]
```

## Description

Restore an existing In-Memory DB Replica Set from one of its point-in-time snapshots.

You pick the source snapshot with --snapshot-id and the replica set to restore it onto with --replicaset-id; that replica set's data is rolled back to the snapshot's state. The snapshot and the replica set must live in the same location/datacenter. Optionally attach a --name and --description to label the restore operation.

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id DisplayName Description ReplicasetId State RestoreTime RestoredSnapshotId]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
      --description string      Optional free-text description of this restore operation
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string             Optional human-readable name to label this restore operation
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
      --replica-set-id string   The ID of the target Replica Set to restore the snapshot onto (must be in the same location as the snapshot) (required)
      --snapshot-id string      The ID of the source Snapshot to restore from (required) (required)
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Restore a replica set from one of its snapshots
ionosctl dbaas in-memory-db snapshot restore create --snapshot-id SNAPSHOT_ID --replicaset-id REPLICASET_ID

# Restore and label the operation
ionosctl dbaas in-memory-db snapshot restore create --snapshot-id SNAPSHOT_ID --replicaset-id REPLICASET_ID --name nightly-rollback --description "roll back after bad deploy"
```

