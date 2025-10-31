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
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ZonesUsed ZonesLimit SecondaryZonesUsed SecondaryZonesLimit RecordsUsed RecordsLimit ReverseRecordsUsed ReverseRecordsLimit]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers        Don't print table headers when table output is used
      --offset int        pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dns quota get
```

