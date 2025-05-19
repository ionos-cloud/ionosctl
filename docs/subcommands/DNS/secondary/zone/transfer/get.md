---
description: "Get the transfer status for a secondary zone"
---

# DnsSecondaryZoneTransferGet

## Usage

```text
ionosctl dns secondary-zone transfer get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get the transfer status for a secondary zone

## Options

```text
  -u, --api-url string    Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [PrimaryIP Status ErrorMessage] (default [PrimaryIP,Status,ErrorMessage])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
  -z, --zone string       The name or ID of the DNS zone
```

