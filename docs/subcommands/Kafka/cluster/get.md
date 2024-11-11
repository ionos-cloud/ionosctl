---
description: "Retrieve a cluster"
---

# KafkaClusterGet

## Usage

```text
ionosctl kafka cluster get [flags]
```

## Aliases

For `cluster` command:

```text
[cl]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a cluster

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The ID of the cluster you want to retrieve (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --location string     The datacenter location (required)
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl kafka cl get --cluster-id ID
```

