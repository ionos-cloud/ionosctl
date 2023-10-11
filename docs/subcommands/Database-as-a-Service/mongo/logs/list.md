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
[mongodb mdb m]
```

For `list` command:

```text
[ls]
```

## Description

List (and optionally filter) the logs of your Mongo Cluster. Use --cols message to see the logs messages.

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Instance Name MessageNumber Message Time]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --direction string    The direction in which to scan through the logs. The logs are returned in order of the direction. Can be one of: BACKWARD, FORWARD
      --end duration        The end time, as a duration. This should be negative and greater than the start time, i.e. -24h. Valid: h, m, s
      --endDate string      The end time for the query in RFC3339 format. Must not be greater than the start parameter. The default value is the current timestamp.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Must be between 1 - 5000 (default 100)
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --start duration      The start time, as a duration. This should be negative, i.e. -720h. Valid: h, m, s
      --startDate string    The start time for the query in RFC3339 format. Must not be greater than 30 days ago and less than the end parameter. The default value is 30 days ago.
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo logs list --cluster-id CLUSTER_ID --start -24h --end -20h --limit 1 --direction FORWARD --cols message
```

