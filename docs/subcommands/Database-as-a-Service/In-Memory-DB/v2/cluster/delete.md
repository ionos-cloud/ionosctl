---
description: "Delete an In-Memory DB Cluster"
---

# DbaasInMemoryDbV2ClusterDelete

## Usage

```text
ionosctl dbaas in-memory-db-v2 cluster delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified In-Memory DB Cluster from your account.

Required values to run command:

* Cluster Id

## Options

```text
  -a, --all                 Delete all Clusters
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'inmemorydbv2' and env var 'IONOS_API_URL' (default "https://in-memory-db.%s.ionos.com")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName DnsName Version Instances Cores Ram State EvictionPolicy PersistenceMode Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string         Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --state string        When used with --all, only delete clusters in this state. Can be one of: PROVISIONING, AVAILABLE, UPDATING, DESTROYING, FAILED
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas in-memory-db-v2 cluster delete --cluster-id <cluster-id>
```

