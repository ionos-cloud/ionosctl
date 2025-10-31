---
description: "Create a kafka topic"
---

# KafkaTopicCreate

## Usage

```text
ionosctl kafka topic create [flags]
```

## Aliases

For `topic` command:

```text
[cl]
```

For `create` command:

```text
[c post]
```

## Description

Create a kafka topic

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --cluster-id string          The ID of the cluster (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name ReplicationFactor NumberOfPartitions RetentionTime SegmentByes State]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --limit int                  Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string            Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
  -n, --name string                The name of the topic (required)
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string              Desired output format [text|json|api-json] (default "text")
      --partitions int32           The number of partitions (default 3)
  -q, --quiet                      Quiet output
      --replication-factor int32   The replication factor (default 3)
      --retention-time int32       The retention time in milliseconds (default 604800000)
      --segment-bytes int32        The segment bytes (default 1073741824)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl kafka topic create --location LOCATION --name my-topic --cluster-id CLUSTER_ID --partitions 1 --replication-factor 1
```

