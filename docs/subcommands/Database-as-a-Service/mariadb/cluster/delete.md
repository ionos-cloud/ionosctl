---
description: "Delete a MariaDB Cluster by ID"
---

# DbaasMariadbClusterDelete

## Usage

```text
ionosctl dbaas mariadb cluster delete [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar ma]
```

For `cluster` command:

```text
[c]
```

For `delete` command:

```text
[del d]
```

## Description

Permanently delete a MariaDB cluster. This destroys the cluster and all of its automatic backups, so point-in-time restore is no longer possible afterwards - this cannot be undone.

Delete a single cluster with --cluster-id, or use --all to remove every cluster (optionally narrowed with --name). With --all, clusters are gathered across every location unless --location pins one.

## Options

```text
  -a, --all                 Delete every MariaDB cluster (across all locations unless --location is set), instead of a single --cluster-id
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster to delete (along with its backups) (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string         With --all, only delete clusters whose display name contains this value (case-insensitive)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mariadb cluster delete --cluster-id <cluster-id>
ionosctl db mar c d --all
ionosctl db mar c d --all --name <name>
```

