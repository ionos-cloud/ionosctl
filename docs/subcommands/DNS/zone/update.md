---
description: "Update a primary DNS zone"
---

# DnsZoneUpdate

## Usage

```text
ionosctl dns zone update [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `update` command:

```text
[u]
```

## Description

Partially update a primary DNS zone. Only the flags you pass change; the rest are preserved (a GET+PUT that simulates PATCH). Identify the zone by name or ID with --zone.

Common use: --enabled=false to take a zone out of service without deleting its records, or --enabled=true to bring it back.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --description string   New free-text note; not served in DNS
      --enabled              Whether the zone is served. true = IONOS answers lookups; false = kept but not resolved (default true)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string          New domain name for the zone, e.g. example.com
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
  -z, --zone string          The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns zone update --zone example.com --description "moved to prod"
ionosctl dns zone update --zone example.com --enabled=false
```

