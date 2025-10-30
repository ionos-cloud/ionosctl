---
description: "Get a MariaDB Cluster by ID"
---

# DbaasMariadbClusterGet

## Usage

```text
ionosctl dbaas mariadb cluster get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get a MariaDB Cluster by ID

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime] (default [ClusterId,Name,DNS,Instances,Version,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci (default "de/txl")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mariadb cluster get --cluster-id <cluster-id>
```

