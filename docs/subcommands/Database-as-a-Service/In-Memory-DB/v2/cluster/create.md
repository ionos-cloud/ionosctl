---
description: "Create an In-Memory DB Cluster"
---

# DbaasInMemoryDbV2ClusterCreate

## Usage

```text
ionosctl dbaas in-memory-db-v2 cluster create [flags]
```

## Aliases

For `in-memory-db-v2` command:

```text
[inmemorydb-v2 memdb-v2 imdb-v2 in-mem-db-v2 inmemdb-v2]
```

For `cluster` command:

```text
[c]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a new In-Memory DB Cluster. The mode is determined by the number of replicas: one replica is standalone, everything else is a replication in leader-follower mode with one active and n-1 passive replicas.

There are two ways to create a cluster, both requiring the same connection and credential flags (--datacenter-id, --lan-id, --cidr, --user, --password) plus --location:
  1. Empty cluster: pass --version (defaults to a supported version) and sizing flags (--replicas, --cores, --ram).
  2. From a snapshot: additionally pass --snapshot-id. The cluster version is taken from the snapshot (so --version is not needed; if given, it must match the snapshot's version). Optionally pass --recovery-time to restore to a point in time within the snapshot's window.

PersistenceMode:
None: Data is inMemory only and will not be persisted. Useful for cache only applications.
AOF (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset.
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
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydbv2' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
  -B, --backup-location string    The Object Storage location where snapshots (backups) will be stored. For added data safety, use a different location than the cluster (default "eu-central-4")
      --cidr string               The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.210.0.0/16 10.212.0.0/14. e.g.: 192.168.1.100/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName DnsName Version Instances Cores Ram State EvictionPolicy PersistenceMode Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username StatusMessage]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance (default 1)
      --datacenter-id string      The unique ID of the Datacenter to connect to your cluster (required)
  -D, --depth int                 Level of detail for response objects (default 1)
      --description string        Human-readable description for the cluster
      --eviction-policy string    The eviction policy for the cluster (refer to the long description for more details). Can be one of: noeviction, allkeys-lru, allkeys-lfu, allkeys-random, volatile-lru, volatile-lfu, volatile-random, volatile-ttl (default "allkeys-lru")
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --lan-id string             The unique ID of the LAN to connect your cluster to (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par
      --logs-enabled              Enable collection and reporting of logs for this cluster
      --maintenance-day string    Day of the week for the MaintenanceWindow. Defaults to a random day during Mon-Fri. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --metrics-enabled           Enable collection and reporting of metrics for this cluster
  -n, --name string               The friendly name of your cluster (default "UnnamedCluster")
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           Password for the initial user. Plaintext is hashed (SHA-256) client-side before sending, as the API only accepts hashed passwords; a value that is already a SHA-256 hash is sent as-is (required)
      --persistence-mode string   Specifies how and if data is persisted (refer to the long description for more details). Can be one of: None, AOF, RDB, RDB_AOF (default "RDB")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance in gigabytes (GB). e.g. --ram 4, --ram 4GB (default "4GB")
      --recovery-time string      Advanced: with --snapshot-id, restore to a specific point in time WITHIN the snapshot's recovery window (PITR), as an ISO 8601 timestamp. Defaults to the latest point in the window
      --replicas int              The total number of replicas in the cluster (one active and n-1 passive). In case of a standalone instance, the value is 1 (default 1)
      --retention-days int32      The number of days snapshots are retained before being automatically deleted (default 7)
      --snapshot-hours ints       Hours of the day (UTC, 0-23) at which snapshots are scheduled to be taken. At least one hour must be specified (default [4])
      --snapshot-id string        Create the cluster from this snapshot instead of empty. The connection/credential flags and --location are still required; the cluster version is taken from the snapshot
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               The initial username (required)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The In-Memory DB version of your Cluster. Ignored when --snapshot-id is set (the snapshot's version is used) (required) (default "8.0")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an empty cluster
ionosctl dbaas in-memory-db-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --version <version>

# Create a cluster from an existing snapshot (version is taken from the snapshot)
ionosctl dbaas in-memory-db-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --snapshot-id <snapshot-id>
```

