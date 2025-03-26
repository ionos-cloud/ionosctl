---
description: "Retrieve your quotas"
---

# DnsQuotaGet

## Usage

```text
ionosctl dns quota get [flags]
```

## Aliases

For `quota` command:

```text
[q]
```

For `get` command:

```text
[g]
```

## Description

Retrieve your quotas

## Options

```text
  -u, --api-url string    Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ZonesUsed ZonesLimit SecondaryZonesUsed SecondaryZonesLimit RecordsUsed RecordsLimit ReverseRecordsUsed ReverseRecordsLimit]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for polling the request (default 60)
  -v, --verbose           Print step-by-step process when running command
  -w, --wait              Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl dns quota get
```

