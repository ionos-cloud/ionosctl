---
description: "Create an In-Memory DB Restore"
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

Create an In-Memory DB Restore

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id ReplicasetId DatacenterId Time State RestoredSnapshotId]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string      A description of the snapshot
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
  -n, --name string             The human readable name of your snapshot
      --no-headers              Don't print table headers when table output is used
      --offset int              pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --replica-set-id string   The ID of the replica set the restore was applied on (required)
      --snapshot-id string      The ID of the In-Memory DB Snapshot to list restore points from (required) (required)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas inmemorydb restore create --replica-set-id REPLICA_SET_ID --snapshot-id SNAPSHOT_ID --name NAME --description DESCRIPTION 
```

