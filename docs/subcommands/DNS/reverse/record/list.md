---
description: "Retrieve all reverse records"
---

# DnsReverseRecordList

## Usage

```text
ionosctl dns reverse-record list [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `list` command:

```text
[ls l]
```

## Description

Retrieve all reverse records

## Options

```text
  -u, --api-url string      Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -i, --ips string          Optional filter for the IP address of the reverse record
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra
      --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        The first element (of the total list of elements) to include in the response. Use together with limit for pagination
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dns rr list
```

