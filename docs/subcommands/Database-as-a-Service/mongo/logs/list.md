---
description: "List (and optionally filter) the logs of your Mongo Cluster. Use --cols message to see the logs messages."
---

# DbaasMongoLogsList

## Usage

```text
ionosctl dbaas mongo logs list [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `list` command:

```text
[ls]
```

## Description

Fetch MongoDB server log lines for a cluster, flattened to one row per message (instance, name, message number, message, time).

Bound the query with a time range given EITHER as absolute RFC3339 timestamps (--startDate/--endDate) OR as relative negative durations from now (--start/--end); the absolute and relative forms of each bound are mutually exclusive. The window may reach at most 30 days into the past. --direction sets scan order (BACKWARD from newest, or FORWARD from oldest) and --limit caps the number of returned lines (1-5000). The message text is hidden by default - add --cols message (or --cols Message) to see it.

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster whose logs to fetch (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Instance Name MessageNumber Message Time]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
      --direction string    Scan order: BACKWARD returns newest-first, FORWARD oldest-first. Determines which end is truncated when --limit is hit. Can be one of: BACKWARD, FORWARD
      --end duration        Relative end of the window as a negative duration from now, e.g. -24h (must be later than the start). Units: h, m, s. Mutually exclusive with --endDate
      --endDate string      Absolute end of the window (RFC3339). Must be after the start. Mutually exclusive with --end. Defaults to now
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of log lines to return (1-5000). When reached, lines are cut at the end according to --direction (default 100)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --start duration      Relative start of the window as a negative duration from now, e.g. -720h. Units: h, m, s. Mutually exclusive with --startDate; window may not start more than 30 days ago
      --startDate string    Absolute start of the window (RFC3339), e.g. 2024-01-15T10:00:00Z. Must be within the last 30 days and before the end. Mutually exclusive with --start. Defaults to 30 days ago
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo logs list --cluster-id CLUSTER_ID --start -24h --end -20h --limit 1 --direction FORWARD --cols message
```

