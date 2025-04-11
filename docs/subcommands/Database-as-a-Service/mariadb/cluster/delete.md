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

Delete a MariaDB Cluster by ID

## Options

```text
  -a, --all                 Delete all mariadb clusters
  -u, --api-url string      Override default host URL (default "https://mariadb.de-txl.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci
      --name                When deleting all clusters, filter the clusters by a name
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl dbaas mariadb cluster delete --cluster-id <cluster-id>
ionosctl db mar c d --all
ionosctl db mar c d --all --name <name>
```

