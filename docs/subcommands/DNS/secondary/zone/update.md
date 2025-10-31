---
description: "Update or create a secondary zone"
---

# DnsSecondaryZoneUpdate

## Usage

```text
ionosctl dns secondary-zone update [flags]
```

## Aliases

For `secondary-zone` command:

```text
[secondary-zones sz]
```

## Description

Update or create a secondary zone

## Options

```text
  -u, --api-url string        Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Id Name Description PrimaryIPs State] (default [Id,Name,Description,PrimaryIPs,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string    Description of the secondary zone
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --limit int             pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string       Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers            Don't print table headers when table output is used
      --offset int            pagination offset: Number of items to skip before starting to collect the results
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --primary-ips strings   Primary DNS server IP addresses
  -q, --quiet                 Quiet output
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string           The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns secondary-zone update --zone ZONE_ID --name ZONE_NAME --description DESCRIPTION --primary-ips 1.2.3.4,5.6.7.8
```

