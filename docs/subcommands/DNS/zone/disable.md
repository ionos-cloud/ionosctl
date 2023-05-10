---
description: Disable a zone
---

# DnsZoneDisable

## Usage

```text
ionosctl dns zone disable [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `disable` command:

```text
[off]
```

## Description

Disable a zone

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -i, --zone-id string   The ID (UUID) of the DNS zone (required)
```

## Examples

```text
ionosctl dns zone disable --zone-id ZONE_ID
```

