---
description: "Create a kafka cluster. Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka"
---

# KafkaClusterCreate

## Usage

```text
ionosctl kafka cluster create [flags]
```

## Aliases

For `cluster` command:

```text
[cl]
```

For `create` command:

```text
[c post]
```

## Description

Create a kafka cluster. Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --broker-addresses strings   The list of broker addresses (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string       The ID of the datacenter (required)
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --lan-id string              The ID of the LAN (required)
      --limit int                  Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string            Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
  -n, --name string                The name of the kafka cluster (required)
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                      Quiet output
      --size string                The size of the kafka cluster. Can be one of: XS, S, M, L, XL (required)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
      --version string             The version of the kafka cluster (required)
```

## Examples

```text
ionosctl kafka cl create --name my-cluster --version 3.7.0 --size S --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 127.0.0.1/24,127.0.0.2/24,127.0.0.3/24
```

