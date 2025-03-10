---
description: "Retrieve a zone"
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
  -u, --api-url string    Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
  -z, --zone string       The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns z get --zone ZONE
```

