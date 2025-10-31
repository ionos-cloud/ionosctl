---
description: "Create a replica set"
---

# DbaasInMemoryDbReplicasetCreate

## Usage

```text
ionosctl dbaas in-memory-db replicaset create [flags]
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

For `create` command:

```text
[post c]
```

## Description

Create a replica set. In-Memory DB replica set with support for a single instance or a In-Memory DB replication in leader follower mode. The mode is determined by the number of replicas. One replica is standalone, everything else an In-Memory DB replication as leader follower mode with one active and n-1 passive replicas.

PersistenceMode:
None: Data is inMemory only and will not be persisted. Useful for cache only applications.
AOF (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the In-Memory DB protocol itself.
RDB: RDB persistence performs snapshots of the current in memory state.
RDB_AOF: Both RDB and AOF persistence are enabled.

EvictionPolicy:
noeviction: No eviction policy is used. In-Memory DB will never remove any data. If the memory limit is reached, an error will be returned on write operations.
allkeys-lru: The least recently used keys will be removed first.
allkeys-lfu: The least frequently used keys will be removed first.
allkeys-random: Random keys will be removed.
volatile-lru: The least recently used keys will be removed first, but only among keys with the expire field set to true.
volatile-lfu: The least frequently used keys will be removed first, but only among keys with the expire field set to true.
volatile-random: Random keys will be removed, but only among keys with the expire field set to true.
volatile-ttl: The key with the nearest time to live will be removed first, but only among keys with the expire field set to true.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --backup-location string    The S3 location where the backups will be stored
      --cidr string               The IP and subnet for your instance. Note the following unavailable IP ranges: 10.210.0.0/16 10.212.0.0/14 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name Version DNSName Replicas Cores RAM StorageSize State BackupLocation PersistenceMode EvictionPolicy MaintenanceDay MaintenanceTime DatacenterId LanId Username]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance (required) (default 1)
      --datacenter-id string      The datacenter to connect your instance to (required)
      --eviction-policy string    The eviction policy for the replica set (refer to the long description for more details). Can be one of: noeviction, allkeys-lru, allkeys-lfu, allkeys-random, volatile-lru, volatile-lfu, volatile-random, volatile-ttl (default "allkeys-lru")
  -f, --force                     Force command to execute without user input
      --hash-password             Hash plaintext passwords before sending. Use '--hash-password=false' to send plaintext passwords as-is (default true)
  -h, --help                      Print usage
      --lan-id string             The numeric Private LAN ID to connect your instance to (required)
      --limit int                 pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/txl, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               The name of the Replica Set (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           Password (plaintext or SHA-256). If plaintext, it’s hashed when --hash-password is true; otherwise sent as-is (required)
      --persistence-mode string   Specifies how and if data is persisted (refer to the long description for more details). Can be one of: None, AOF, RDB, RDB_AOF (default "RDB")
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance in gigabytes (GB) (required) (default "4GB")
      --replicas int              The total number of replicas in the Replica Set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is >1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance (required) (default 1)
      --snapshot-id string        If set, will create the replicaset from the specified snapshot
      --user string               The initial username (required)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The In-Memory DB version of your Replica Set (required) (default "7.2")
```

## Examples

```text
ionosctl dbaas inmemorydb replicaset create --location LOCATION --name NAME --replicas REPLICAS --cores CORES --ram RAM --user USER --password PASSWORD --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 
```

