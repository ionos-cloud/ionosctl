---
description: "Create a primary DNS zone"
---

# DnsZoneCreate

## Usage

```text
ionosctl dns zone create [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `create` command:

```text
[post c]
```

## Description

Create a primary DNS zone for a domain you want IONOS to answer for.

--name is the domain itself (e.g. example.com), NOT a friendly label. After creating the zone, delegate the domain to the IONOS name servers at your registrar and add entries with 'dns record create'. A zone starts --enabled; pass --enabled=false to create it dormant.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --description string   Free-text note for your own reference; not served in DNS
      --enabled              Whether the zone is served. true = IONOS answers lookups; false = zone kept but not resolved (default true) (default true)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string          Domain name this zone is authoritative for, e.g. example.com (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dns zone create --name example.com
ionosctl dns zone create --name example.com --description "prod apex" --enabled=false
```

