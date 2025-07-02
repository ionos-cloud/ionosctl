---
description: "List all kafka topics"
---

# KafkaTopicList

## Usage

```text
ionosctl kafka topic list [flags]
```

## Aliases

For `topic` command:

```text
[cl]
```

For `list` command:

```text
[ls]
```

## Description

List all kafka topics

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --cluster-id string   The ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name ReplicationFactor NumberOfPartitions RetentionTime SegmentByes State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl (default "de/fra")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl kafka topic list --location LOCATION
ionosctl kafka topic list --location LOCATION --cluster-id CLUSTER_ID
```

