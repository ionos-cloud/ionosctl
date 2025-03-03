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
  -u, --api-url string        Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Id Name Description PrimaryIPs State] (default [Id,Name,Description,PrimaryIPs,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --description string    Description of the secondary zone
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -l, --location string       Location of the resource to operate on. Can be one of: de/fra
      --no-headers            Don't print table headers when table output is used
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --primary-ips strings   Primary DNS server IP addresses
  -q, --quiet                 Quiet output
  -v, --verbose count         Print step-by-step process when running command
  -z, --zone string           The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns secondary-zone update --zone ZONE_ID --name ZONE_NAME --description DESCRIPTION --primary-ips 1.2.3.4,5.6.7.8
```

