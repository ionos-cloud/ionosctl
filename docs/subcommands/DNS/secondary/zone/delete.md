---
description: "Delete a secondary zone"
---

# DnsSecondaryZoneDelete

## Usage

```text
ionosctl dns secondary-zone delete [flags]
```

## Aliases

For `secondary-zone` command:

```text
[secondary-zones sz]
```

For `delete` command:

```text
[d del]
```

## Description

Delete a secondary zone

## Options

```text
  -a, --all               Delete all secondary zones
  -u, --api-url string    Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://dns.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Description PrimaryIPs State] (default [Id,Name,Description,PrimaryIPs,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
  -z, --zone string       The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns secondary-zone delete --zone ZONE_ID
```

