---
description: "Create a secondary DNS zone"
---

# DnsSecondaryZoneCreate

## Usage

```text
ionosctl dns secondary-zone create [flags]
```

## Aliases

For `secondary-zone` command:

```text
[secondary-zones sz]
```

For `create` command:

```text
[c]
```

## Description

Create a secondary zone: a read-only copy that IONOS pulls from your external primary name server.

--name is the domain (e.g. example.com); --primary-ips lists the primary servers IONOS transfers from. The zone starts with default NS/SOA records; run 'dns secondary-zone transfer start' to pull the rest.

IONOS CLOUD DNS sends its DNS NOTIFY messages from these Anycast addresses — allow them on your primary:
  IPv4: 212.227.123.25
  IPv6: 2001:8d8:fe:53::5cd:25

## Options

```text
  -u, --api-url string        Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Id Name Description PrimaryIPs State]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int             Level of detail for response objects (default 1)
      --description string    Free-text note for your own reference; not served in DNS
  -F, --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --limit int             Maximum number of items to return per request (default 50)
  -l, --location string       Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string           Domain name of the zone to mirror, e.g. example.com
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --primary-ips strings   Comma-separated IPs of the external primary name servers IONOS transfers the zone from
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
  -t, --timeout int           Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                  Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dns secondary-zone create --name example.com --primary-ips 1.2.3.4,5.6.7.8
```

