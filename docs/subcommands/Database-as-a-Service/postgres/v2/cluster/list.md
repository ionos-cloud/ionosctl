---
description: "List PostgreSQL Clusters"
---

# DbaasPostgresV2ClusterList

## Usage

```text
ionosctl dbaas postgres-v2 cluster list [flags]
```

## Aliases

For `postgres-v2` command:

```text
[pg-v2 pgsql-v2 postgresql-v2 psql-v2]
```

For `cluster` command:

```text
[c]
```

For `list` command:

```text
[ls]
```

## Description

Use this command to retrieve a list of PostgreSQL Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'psqlv2' and env var 'IONOS_API_URL' (default "https://postgresql.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ClusterId DisplayName DnsName PostgresVersion Instances Ram Cores StorageSize State SyncMode Description ConnectionPooler MaintenanceDay MaintenanceTime BackupLocation LogsEnabled MetricsEnabled DatacenterId LanId Cidr DbUsername DbDatabase StatusMessage] (default [ClusterId,DisplayName,DnsName,PostgresVersion,Instances,Ram,Cores,StorageSize,State,SyncMode])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int32       The maximum number of elements to return (default 100)
  -l, --location string   Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, gb/bhx, us/las, us/mci, us/ewr (default "de/txl")
  -n, --name string       Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers        Don't print table headers when table output is used
      --offset int32      The first element to return
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --state string      Response filter by cluster state: PROVISIONING, AVAILABLE, UPDATING, DESTROYING, FAILED
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres-v2 cluster list
```

