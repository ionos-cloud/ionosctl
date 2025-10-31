---
description: "Retrieve all records from either a primary or secondary zone"
---

# DnsRecordList

## Usage

```text
ionosctl dns record list [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve all records from either a primary or secondary zone

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns for primary zones: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
                                Available columns for secondary zones: [Id Name Content Type Enabled FQDN ZoneId ZoneName RootName]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string             Filter used to fetch only the records that contain specified record name. NOTE: Only available for zone records.
      --no-headers              Don't print table headers when table output is used
      --offset int              Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --secondary-zone string   The name or ID of the secondary zone to fetch records from
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string             (UUID or Zone Name) Filter used to fetch only the records that contain specified zone.
```

## Examples

```text
ionosctl dns r list
ionosctl dns r list --secondary-zone SECONDARY_ZONE_ID
ionosctl dns r list --zone ZONE_ID
```

