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
  -u, --api-url string          Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns for primary zones: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
                                Available columns for secondary zones: [Id Name Content Type Enabled FQDN ZoneId ZoneName RootName]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra
      --max-results int32       The maximum number of elements to return
      --name string             Filter used to fetch only the records that contain specified record name. NOTE: Only available for zone records.
      --no-headers              Don't print table headers when table output is used
      --offset int32            The first element (of the total list of elements) to include in the response. Use together with limit for pagination
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --secondary-zone string   The name or ID of the secondary zone to fetch records from
  -t, --timeout duration        Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait                    Polls the request continuously until the operation is completed
  -z, --zone string             (UUID or Zone Name) Filter used to fetch only the records that contain specified zone.
```

## Examples

```text
ionosctl dns r list
ionosctl dns r list --secondary-zone SECONDARY_ZONE_ID
ionosctl dns r list --zone ZONE_ID
```

