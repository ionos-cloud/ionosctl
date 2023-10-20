---
description: "Delete a zone"
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
  -a, --all              Delete all zones. Required or -z
  -u, --api-url string   Override default host url (default "dns.de-fra.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -z, --zone string      The name or ID of the DNS zone. Required or -a
```

## Examples

```text
ionosctl dns z delete --zone ZONE
```

