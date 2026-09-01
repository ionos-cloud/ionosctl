---
description: "Create a Kafka topic"
---

# KafkaTopicCreate

## Usage

```text
ionosctl kafka topic create [flags]
```

## Aliases

For `topic` command:

```text
[t]
```

For `create` command:

```text
[c post]
```

## Description

Create a topic inside a Kafka cluster.

--partitions splits the topic so it can be produced/consumed in parallel; ordering is guaranteed only within a partition, and partitions cannot be reduced later.

--replication-factor is how many brokers keep a copy of each partition, giving fault tolerance. It cannot exceed the number of brokers in the cluster (3), so valid values are 1-3; 3 is recommended for production.

Log retention decides how long data is kept: --retention-time is the age (in milliseconds) after which messages become eligible for deletion (default 604800000 = 7 days); --segment-bytes is the size each on-disk log segment file reaches before a new one is rolled (default 1073741824 = 1 GiB). Retained data counts against the cluster's shared storage.

The cluster must be AVAILABLE before topics can be created.

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --cluster-id string          ID of the cluster to create the topic in (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name ReplicationFactor NumberOfPartitions RetentionTime SegmentByes ClusterId State]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                  Level of detail for response objects (default 1)
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --limit int                  Maximum number of items to return per request (default 50)
  -l, --location string            Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string                Name of the topic (required)
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
      --partitions int32           Number of partitions the topic is split into (parallelism / ordering unit); cannot be reduced later (default 3)
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --replication-factor int32   Copies of each partition kept across brokers for fault tolerance; 1-3 (cannot exceed the 3 brokers) (default 3)
      --retention-time int32       Age in milliseconds after which messages may be deleted (default 604800000 = 7 days) (default 604800000)
      --segment-bytes int32        Size in bytes a log segment file reaches before a new one is rolled (default 1073741824 = 1 GiB) (default 1073741824)
  -t, --timeout int                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl kafka topic create --location LOCATION --cluster-id CLUSTER_ID --name my-topic --partitions 3 --replication-factor 3
ionosctl kafka topic create --location LOCATION --cluster-id CLUSTER_ID --name events --partitions 6 --replication-factor 3 --retention-time 86400000 --segment-bytes 536870912
```

