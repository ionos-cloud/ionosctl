---
description: "Create a secondary zone"
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

Create a new secondary zone with default NS and SOA records. Note that Cloud DNS relies on the following Anycast addresses for sending DNS notify messages. Make sure to whitelist on your end:

IPv4: 212.227.123.25
IPv6: 2001:8d8:fe:53::5cd:25

## Options

```text
  -u, --api-url string        Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Id Name Description PrimaryIPs State] (default [Id,Name,Description,PrimaryIPs,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int             Level of detail for response objects (default 1)
      --description string    Description of the secondary zone
      --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --limit int             Maximum number of items to return per request (default 50)
  -l, --location string       Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string           Name of the secondary zone
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --primary-ips strings   Primary DNS server IP addresses
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dns secondary-zone create --name ZONE_NAME --description DESCRIPTION --primary-ips 1.2.3.4,5.6.7.8
```

