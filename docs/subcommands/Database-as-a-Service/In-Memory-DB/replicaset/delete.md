---
description: "Delete In-Memory DB Replica Sets"
---

# DbaasInMemoryDbReplicasetDelete

## Usage

```text
ionosctl dbaas in-memory-db replicaset delete [flags]
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

For `delete` command:

```text
[d del remove rm]
```

## Description

Delete In-Memory DB Replica Sets

## Options

```text
  -a, --all                     Delete all replica-sets. Required or -replica-set-id
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Name Version DNSName Replicas Cores RAM StorageSize State BackupLocation PersistenceMode EvictionPolicy MaintenanceDay MaintenanceTime DatacenterId LanId Username]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -i, --replica-set-id string   The ID of the Replica Set you want to delete
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas inmemorydb replicaset delete --replica-set-id REPLICA_SET_ID --force FORCE 
ionosctl dbaas inmemorydb replicaset delete --all ALL --force FORCE 
```

