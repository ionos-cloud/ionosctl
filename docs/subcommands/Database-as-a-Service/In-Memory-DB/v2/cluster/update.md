---
description: "Update an In-Memory DB Cluster"
---

# DbaasInMemoryDbV2ClusterUpdate

## Usage

```text
ionosctl dbaas in-memory-db-v2 cluster update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update attributes of an In-Memory DB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydbv2' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId DisplayName DnsName Version Instances Cores Ram State EvictionPolicy PersistenceMode Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username StatusMessage]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The number of CPU cores per instance
  -D, --depth int                 Level of detail for response objects (default 1)
      --description string        Human-readable description for the cluster
      --eviction-policy string    The eviction policy for the cluster. Can be one of: noeviction, allkeys-lru, allkeys-lfu, allkeys-random, volatile-lru, volatile-lfu, volatile-random, volatile-ttl
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --logs-enabled              Enable collection and reporting of logs for this cluster
      --maintenance-day string    Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
      --maintenance-time string   Time for the MaintenanceWindow. e.g.: 16:30:59. Must be specified together with --maintenance-day
      --metrics-enabled           Enable collection and reporting of metrics for this cluster
  -n, --name string               The friendly name of your cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           Password for the In-Memory DB user. Required because the API does not return it on GET requests. Plaintext is hashed (SHA-256) client-side (required)
      --persistence-mode string   Specifies how and if data is persisted. Can be one of: None, AOF, RDB, RDB_AOF
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                The amount of memory per instance in GB. e.g. --ram 4, --ram 4GB
      --replicas int              The total number of replicas in the cluster (one active and n-1 passive)
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               Username for the In-Memory DB user. Defaults to the cluster's current username
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The In-Memory DB version of your cluster
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas in-memory-db-v2 cluster update --cluster-id <cluster-id> --password <password> --cores 4 --ram 8GB
```

