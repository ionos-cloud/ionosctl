---
description: "Retrieve all clusters using pagination and optional filters"
---

# KafkaClusterList

## Usage

```text
ionosctl kafka cluster list [flags]
```

## Aliases

For `cluster` command:

```text
[cl]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve all clusters using pagination and optional filters

## Options

```text
  -u, --api-url string    Override default host URL (default "https://kafka.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra, de/txl
      --name string       Filter used to fetch only the records that contain specified name.
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
      --state string      Filter used to fetch only the records that contain specified state.. Can be one of: AVAILABLE, BUSY, DEPLOYING, UPDATING, FAILED_UPDATING, FAILED, DESTROYING
  -v, --verbose count     Print step-by-step process when running command
```

## Examples

```text
ionosctl kafka c list --location de/txl
```

