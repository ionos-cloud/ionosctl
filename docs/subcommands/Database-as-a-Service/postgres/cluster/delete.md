---
description: "Delete a PostgreSQL Cluster"
---

# DbaasPostgresClusterDelete

## Usage

```text
ionosctl dbaas postgres cluster delete [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
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

Use this command to delete a specified PostgreSQL Cluster from your account. You can wait for the cluster to be deleted with the wait-for-deletion option.

Required values to run command:

* Cluster Id

## Options

```text
  -a, --all                 Delete all Clusters
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName Location State PostgresVersion Instances RAM Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -n, --name string         Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for Cluster to be completely removed[seconds] (default 1200)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -W, --wait-for-deletion   Wait for Cluster to be completely removed
```

## Examples

```text
ionosctl dbaas postgres cluster delete -i CLUSTER_ID
```

