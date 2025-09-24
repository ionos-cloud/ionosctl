---
description: "List In-Memory DB Snapshots"
---

# DbaasInMemoryDbSnapshotList

## Usage

```text
ionosctl dbaas in-memory-db snapshot list [flags]
```

## Aliases

For `in-memory-db` command:

```text
[inmemorydb memdb imdb in-mem-db inmemdb]
```

For `snapshot` command:

```text
[snaps snap backup backups snapshots]
```

For `list` command:

```text
[l ls]
```

## Description

List In-Memory DB Snapshots

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id ReplicasetId DatacenterId Time State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas inmemorydb snapshot list
```

