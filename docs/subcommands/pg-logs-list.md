---
description: List Logs for a PostgreSQL Cluster
---

# PgLogsList

## Usage

```text
ionosctl pg logs list [flags]
```

## Aliases

For `pg` command:

```text
[postgres]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve the Logs of a specified PostgreSQL Cluster. By default, the result will contain all Cluster Logs. You can specify the start time, end time or a limit for sorting Cluster Logs.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Name Message Time Logs] (default [Logs])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -e, --end-time string     The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --limit int           The maximal number of log lines to return. The command will print all logs, if this is not set
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -s, --start-time string   The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl pg logs list --cluster-id CLUSTER_ID
```

