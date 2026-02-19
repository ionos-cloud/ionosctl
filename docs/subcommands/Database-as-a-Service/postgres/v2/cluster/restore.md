---
description: "Restore a PostgreSQL Cluster"
---

# DbaasPostgresV2ClusterRestore

## Usage

```text
ionosctl dbaas postgres-v2 cluster restore [flags]
```

## Aliases

For `postgres-v2` command:

```text
[pg-v2]
```

For `cluster` command:

```text
[c]
```

For `restore` command:

```text
[r]
```

## Description

Use this command to trigger an in-place restore of the specified PostgreSQL Cluster.

Required values to run command:

* Cluster Id
* Backup Id

## Options

```text
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'todo' and env var 'IONOS_API_URL' (default "https://postgresql.%s.ionos.com")
      --backup-id string       The unique ID of the backup you want to restore (required)
  -i, --cluster-id string      The unique ID of the Cluster (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId DisplayName DnsName PostgresVersion Instances Ram Cores StorageSize State SyncMode MaintenanceDay MaintenanceTime BackupLocation DatacenterId LanId Cidr] (default [ClusterId,DisplayName,DnsName,PostgresVersion,Instances,Ram,Cores,StorageSize,State,SyncMode])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, gb/bhx, us/las, us/mci, us/ewr (default "de/txl")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -R, --recovery-time string   If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely
  -t, --timeout int            Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -W, --wait-for-state         Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres-v2 cluster restore --cluster-id <cluster-id> --backup-id <backup-id>
```

