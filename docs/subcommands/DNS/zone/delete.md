---
description: Delete a zone
---

# DnsZoneDelete

## Usage

```text
ionosctl dns zone delete [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `delete` command:

```text
[del d]
```

## Description

Delete a zone

## Options

```text
  -a, --all              Delete all zones. Required or -i
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Skip yes/no confirmation
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -i, --zone-id string   The ID (UUID) of the DNS zone. Required or -a
```

## Examples

```text
ionosctl dns zone delete --zone-id ZONE_ID
```

