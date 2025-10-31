---
description: "Get an In-Memory DB Replica Set"
---

# DbaasInMemoryDbReplicasetGet

## Usage

```text
ionosctl dbaas in-memory-db replicaset get [flags]
```

## Aliases

For `in-memory-db` command:

```text
[inmemorydb memdb imdb in-mem-db inmemdb]
```

For `replicaset` command:

```text
[rs replica-set replicasets cluster]
```

For `get` command:

```text
[g]
```

## Description

Get an In-Memory DB Replica Set

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Name Version DNSName Replicas Cores RAM StorageSize State BackupLocation PersistenceMode EvictionPolicy MaintenanceDay MaintenanceTime DatacenterId LanId Username]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers              Don't print table headers when table output is used
      --offset int              Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -i, --replica-set-id string   The ID of the Replica Set you want to delete
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas inmemorydb replicaset get --replica-set-id REPLICA_SET_ID 
```

