---
description: Retrieve a record
---

# DnsRecordGet

## Usage

```text
ionosctl dns record get [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a record

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [Id Name Content Type Enabled FQDN State ZoneId ZoneName]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         When using text output, don't print headers
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --record-id string   The ID (UUID) of the DNS record
  -v, --verbose            Print step-by-step process when running command
  -z, --zone string        The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns record get --zone ZONE --recordId RECORD_ID
```

