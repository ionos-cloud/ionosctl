---
description: "List MariaDB Backups"
---

# DbaasMariadbBackupList

## Usage

```text
ionosctl dbaas mariadb backup list [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar ma]
```

For `backup` command:

```text
[b]
```

For `list` command:

```text
[l ls]
```

## Description

List all MariaDB Backups, or optionally provide a Cluster ID to list those of a certain cluster

## Options

```text
  -u, --api-url string      Override default host URL (default "https://mariadb.de-txl.ionos.com")
  -i, --cluster-id string   Optionally limit shown backups to those of a certain cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [BackupId ClusterId Size Items] (default [BackupId,ClusterId,Size,Items])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl dbaas mariadb backup list
```

