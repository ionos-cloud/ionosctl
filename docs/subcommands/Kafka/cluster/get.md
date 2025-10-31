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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
  -i, --cluster-id string   The ID of the cluster you want to retrieve (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl kafka cl get --cluster-id ID
```

