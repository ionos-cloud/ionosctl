---
description: "Initiate zone transfer"
---

# DnsSecondaryZoneTransferStart

## Usage

```text
ionosctl dns secondary-zone transfer start [flags]
```

## Aliases

For `secondary-zone` command:

```text
[secondary-zones sz]
```

For `transfer` command:

```text
[t]
```

For `start` command:

```text
[s]
```

## Description

Initiate zone transfer

## Options

```text
  -u, --api-url string   Override default host url (default "dns.de-fra.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [PrimaryIP Status ErrorMessage] (default [PrimaryIP,Status,ErrorMessage])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -z, --zone string      The name or ID of the DNS zone
```
