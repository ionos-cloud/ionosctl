---
description: "List Logs for a PostgreSQL Cluster"
---

# DbaasPostgresLogsList

## Usage

```text
ionosctl dbaas postgres logs list [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
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
  -D, --direction string    The direction in which to scan through the logs. The logs are returned in order of the direction. (default "BACKWARD")
  -e, --end-time string     The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --limit int           The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Minimum: 1. Maximum: 5000 (default 100)
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -S, --since string        The start time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both start-time and since are set, start-time will be used.
  -s, --start-time string   The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z
  -U, --until string        The end time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both end-time and until are set, end-time will be used.
  -v, --verbose count       Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres logs list --cluster-id CLUSTER_ID --since 5h --until 1h
```

