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
  -u, --api-url string             Override default host URL (default "https://kafka.de-fra.ionos.com")
      --broker-addresses strings   The list of broker addresses (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string       The ID of the datacenter (required)
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --lan-id string              The ID of the LAN (required)
  -l, --location string            Location of the resource to operate on. Can be one of: de/fra, de/txl
      --name string                The name of the kafka cluster (required)
      --no-headers                 Don't print table headers when table output is used
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                      Quiet output
      --size string                The size of the kafka cluster. Can be one of: XS, S, M, L, XL (required)
  -t, --timeout int                Timeout in seconds for polling the request (default 60)
  -v, --verbose                    Print step-by-step process when running command
      --version string             The version of the kafka cluster (required)
  -w, --wait                       Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl kafka cl create --name my-cluster --version 3.7.0 --size S --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 127.0.0.1/24,127.0.0.2/24,127.0.0.3/24
```

