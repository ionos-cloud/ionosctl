---
description: "Create an In-Memory DB Replica Set"
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

Create a new In-Memory DB Replica Set (a Redis-compatible database).

The mode is set by --replicas: 1 replica is a standalone instance; more than 1 is a leader-follower replication with one active and n-1 passive standby replicas (passive replicas fail over but do not serve reads). Up to 5 replicas are allowed.

Each replica gets its own --cores (1-31) and --ram (4-256 GB); storage is derived automatically from the RAM and persistence mode and cannot be set. The replica set attaches to exactly one network connection: --datacenter-id + --lan-id + --cidr. An initial user (--user / --password) is created for you.

There are two ways to create a replica set:
  1. Empty: pass the sizing, connection and credential flags (as in the basic example below).
  2. From a snapshot: additionally pass --snapshot-id to restore an existing point-in-time snapshot into the new replica set. The snapshot must belong to the same location.

PersistenceMode (--persistence-mode, controls how data survives restarts):
  None:    In-memory only, nothing is persisted. Best for pure caches.
  AOF:     Append Only File - every write is logged and replayed on restart, reconstructing the dataset.
  RDB:     Periodic point-in-time dumps of the in-memory state.
  RDB_AOF: Both RDB and AOF are enabled.

EvictionPolicy (--eviction-policy, what happens when the memory limit is hit):
  noeviction:      Never evict; write operations return an error once memory is full.
  allkeys-lru:     Evict the least recently used keys first.
  allkeys-lfu:     Evict the least frequently used keys first.
  allkeys-random:  Evict random keys.
  volatile-lru:    As allkeys-lru, but only among keys with a TTL (expire) set.
  volatile-lfu:    As allkeys-lfu, but only among keys with a TTL (expire) set.
  volatile-random: Evict random keys, but only among keys with a TTL (expire) set.
  volatile-ttl:    Evict the key with the nearest TTL first, among keys with a TTL set.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydb' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
      --backup-location string    The S3 (Object Storage) location where automatic backups/snapshots are stored, e.g. 'de'. Defaults to a location near the cluster
      --cidr string               The private IP and subnet assigned to the replica set on the LAN, in CIDR notation (e.g. 192.168.1.100/24). These ranges are reserved and cannot be used: 10.210.0.0/16, 10.212.0.0/14 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name Version DNSName Replicas Cores RAM StorageSize State BackupLocation PersistenceMode EvictionPolicy MaintenanceDay MaintenanceTime DatacenterId LanId Username]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance (1-31) (required) (default 1)
      --datacenter-id string      The ID of the Virtual Datacenter the replica set connects into. Must be in the same location as the replica set (required)
  -D, --depth int                 Level of detail for response objects (default 1)
      --eviction-policy string    What to evict when the memory limit is reached. 'volatile-*' policies only touch keys with a TTL; 'noeviction' errors on writes when full. See the long description for details. Can be one of: noeviction, allkeys-lru, allkeys-lfu, allkeys-random, volatile-lru, volatile-lfu, volatile-random, volatile-ttl (default "allkeys-lru")
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
      --hash-password             Hash a plaintext --password (SHA-256) before sending. Set --hash-password=false to send the plaintext password to the API as-is (default true)
  -h, --help                      Print usage
      --lan-id string             The numeric ID of a private LAN (within the datacenter) to attach the replica set to (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --maintenance-day string    Day of the week for the weekly 4-hour maintenance window (Monday-Sunday). Defaults to a random weekday (Mon-Fri) (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Start time (UTC, HH:MM:SS) of the weekly 4-hour maintenance window during which upgrades and patches may occur, e.g. 16:30:59. Defaults to a random time between 10:00 and 16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               The human-readable display name of the Replica Set (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           Password for the initial user. Plaintext must be 10-63 chars; by default it is hashed (SHA-256) client-side before sending. A value that is already a 64-char SHA-256 hash is sent as-is (required)
      --persistence-mode string   How data is persisted across restarts: None (cache only), AOF (write log), RDB (periodic dumps), or RDB_AOF (both). See the long description for details. Can be one of: None, AOF, RDB, RDB_AOF (default "RDB")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance, 4-256 GB (e.g. --ram 8 or --ram 8GB). Storage size is derived automatically from RAM and persistence mode and is not configurable (required) (default "4GB")
      --replicas int              Total number of replicas (1-5): one active plus n-1 passive. Set 1 for a standalone instance; >1 enables leader-follower replication. Passive replicas are hot standbys for failover of the active instance only - they are NOT read replicas and do not serve reads (required) (default 1)
      --snapshot-id string        Create the replica set restored from this existing snapshot instead of empty. The snapshot must be in the same location; all sizing/connection/credential flags still apply
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               Username for the initial database user. 1-16 chars, letters/digits/underscore only. Some names are reserved (e.g. 'admin', 'standby') (required)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The In-Memory DB (Redis-compatible) engine version. Supported values: 7.0, 7.2 (required) (default "7.2")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a standalone (single-instance) replica set
ionosctl dbaas in-memory-db replicaset create --location de/fra --name mycache --replicas 1 --cores 1 --ram 4GB --datacenter-id DATACENTER_ID --lan-id 1 --cidr 192.168.1.100/24 --user dbadmin --password MyStrongPass1

# Advanced: a 3-node leader-follower set with AOF persistence, a custom eviction policy and maintenance window
ionosctl dbaas in-memory-db replicaset create --location de/fra --name prod-cache --replicas 3 --cores 4 --ram 16GB --persistence-mode RDB_AOF --eviction-policy volatile-lru --datacenter-id DATACENTER_ID --lan-id 1 --cidr 192.168.1.100/24 --user dbadmin --password MyStrongPass1 --maintenance-day Saturday --maintenance-time 03:00:00
```

