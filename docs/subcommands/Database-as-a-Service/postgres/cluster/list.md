---
description: "List PostgreSQL clusters"
---

# DbaasPostgresClusterList

## Usage

```text
ionosctl dbaas postgres cluster list [flags]
```

## Aliases

For `postgres` command:

```text
[pg pgsql postgresql psql]
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

Retrieve all PostgreSQL clusters provisioned under your account. Use `--name` to keep only clusters whose display name contains the given substring (case-insensitive).

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ClusterId DisplayName Location DatacenterId LanId Cidr Instances State PostgresVersion RAM Cores StorageSize StorageType MaintenanceWindow SynchronizationMode BackupLocation]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -n, --name string       Keep only clusters whose display name contains this substring (case-insensitive). Omit to list all clusters
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas postgres cluster list
```

