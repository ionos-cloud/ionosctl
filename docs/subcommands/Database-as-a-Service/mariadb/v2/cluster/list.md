---
description: "List MariaDB Clusters"
---

# DbaasMariadbV2ClusterList

## Usage

```text
ionosctl dbaas mariadb-v2 cluster list [flags]
```

## Aliases

For `mariadb-v2` command:

```text
[maria-v2 mar-v2 ma-v2]
```

For `cluster` command:

```text
[c]
```

For `list` command:

```text
[ls l]
```

## Description

Use this command to retrieve a list of MariaDB Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadbv2' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ClusterId Name DnsName Version Instances State Cores Ram StorageSize Description BackupLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr StatusMessage]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int32       The maximum number of elements to return (default 100)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string       Response filter to list only the MariaDB Clusters that contain the specified name in the Name field. The value is case insensitive
      --no-headers        Don't print table headers when table output is used
      --offset int32      The first element to return
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --state string      Response filter by cluster state: PROVISIONING, AVAILABLE, UPDATING, DESTROYING, FAILED
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mariadb-v2 cluster list
```

