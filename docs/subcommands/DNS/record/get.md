---
description: "Retrieve a record"
---

# DnsRecordGet

## Usage

```text
ionosctl dns record get [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a record

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers        Don't print table headers when table output is used
      --offset int        pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
      --record string     The ID or name of the DNS record
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string       The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns r get --zone ZONE --record RECORD
```

