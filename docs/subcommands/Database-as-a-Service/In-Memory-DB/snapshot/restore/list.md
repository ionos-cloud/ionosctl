---
description: "List In-Memory DB Restores"
---

# DbaasInMemoryDbSnapshotRestoreList

## Usage

```text
ionosctl dbaas in-memory-db snapshot restore list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List In-Memory DB Restores

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'in-memory-db' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id ReplicasetId DatacenterId Time State RestoredSnapshotId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --snapshot-id string   The ID of the In-Memory DB Snapshot to list restore points from (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas inmemorydb restore list
```

