---
description: Retrieve a zone
---

# DnsZoneGet

## Usage

```text
ionosctl dns zone get [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a zone

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -i, --zone-id string   The ID (UUID) of the DNS zone (required)
```

## Examples

```text
ionosctl dns zone get --zone-id ZONE_ID
```

