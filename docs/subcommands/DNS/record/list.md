---
description: "Retrieve all records"
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

Retrieve all records

## Options

```text
  -u, --api-url string      Override default host url (default "dns.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Content Type Enabled FQDN State ZoneId ZoneName]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --max-results int32   The maximum number of elements to return
      --name string         Filter used to fetch only the records that contain specified record name
      --no-headers          Don't print table headers when table output is used
      --offset int32        The first element (of the total list of elements) to include in the response. Use together with limit for pagination
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
  -z, --zone string         (UUID or Zone Name) Filter used to fetch only the records that contain specified zone.
```

## Examples

```text
ionosctl dns r list
```

