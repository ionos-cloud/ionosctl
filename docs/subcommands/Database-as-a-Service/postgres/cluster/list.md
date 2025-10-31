---
description: "List PostgreSQL Clusters"
---

# DbaasPostgresClusterList

## Usage

```text
ionosctl dbaas postgres cluster list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of PostgreSQL Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [ClusterId DisplayName Location State PostgresVersion Instances RAM Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string      Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres cluster list
```

