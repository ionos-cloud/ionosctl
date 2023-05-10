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
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --record-id string   The ID (UUID) of the DNS record
  -v, --verbose            Print step-by-step process when running command
      --zone-id string     The ID (UUID) of the DNS zone of which record belongs to
```

## Examples

```text
ionosctl dns record get --zoneId ZONE_ID --recordId RECORD_ID
```

