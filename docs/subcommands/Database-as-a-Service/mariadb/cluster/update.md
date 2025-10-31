---
description: "Update a MariaDB Cluster"
---

# DbaasMariadbClusterUpdate

## Usage

```text
ionosctl dbaas mariadb cluster update [flags]
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

For `update` command:

```text
[u]
```

## Description

Update a MariaDB Cluster

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
  -i, --cluster-id string         The unique ID of the cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime] (default [ClusterId,Name,DNS,Instances,Version,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               Core count. Can be increased or decreased.
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries). Instances can only be increased (3,5,7)
      --limit int                 pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci (default "de/txl")
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. e.g.: Monday. To change maintenance provide both --maintenance-day and --maintenance-time
      --maintenance-time string   Time for the MaintenanceWindows. e.g.: 16:30:59. To change maintenance provide both --maintenance-day and --maintenance-time
  -n, --name string               The name of your cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                RAM size. e.g.: --ram 4GB. Can be increased or decreased.
      --storage-size string       The size of the Storage in GB. Can only be increased
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The MariaDB version of your cluster. Downgrades are not supported (version can only be increased) 
```

## Examples

```text
ionosctl dbaas mariadb cluster update--Cluster ID: %v CLUSTER ID: %V --version VERSION 
```

