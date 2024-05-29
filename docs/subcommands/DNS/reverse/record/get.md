---
description: "Find a record by IP or ID"
---

# DnsReverseRecordGet

## Usage

```text
ionosctl dns reverse-record get [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `get` command:

```text
[g]
```

## Description

Find a record by IP or ID

## Options

```text
  -u, --api-url string   Override default host url (default "dns.de-fra.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
      --record string    The record ID or IP, for identifying which record you want to update (required)
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dns rr get --record RECORD_IP
ionosctl dns rr get --record RECORD_ID
```

