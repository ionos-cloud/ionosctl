---
description: "Delete a cluster"
---

# KafkaClusterDelete

## Usage

```text
ionosctl kafka cluster delete [flags]
```

## Aliases

For `cluster` command:

```text
[cl]
```

For `delete` command:

```text
[del d]
```

## Description

Delete a cluster

## Options

```text
  -a, --all                 Delete all records if set (required)
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
  -i, --cluster-id string   The ID of the cluster you want to retrieve (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par (default "de/fra")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl kafka cl delete --cluster-id ID
```

