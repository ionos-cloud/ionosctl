---
description: "Get a kafka topic"
---

# KafkaTopicGet

## Usage

```text
ionosctl kafka topic get [flags]
```

## Aliases

For `topic` command:

```text
[cl]
```

For `get` command:

```text
[g]
```

## Description

Get a kafka topic

## Options

```text
  -u, --api-url string      Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://kafka.de-fra.ionos.com")
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
      --topic-id string     The ID of the topic (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl kafka topic get --location LOCATION --cluster-id CLUSTER_ID --topic-id TOPIC_ID
```

