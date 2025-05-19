---
description: "Retrieve zones"
---

# DnsZoneList

## Usage

```text
ionosctl dns zone list [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve zones

## Options

```text
  -u, --api-url string      Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra
      --max-results int32   Pagination limit
      --name string         Filter used to fetch only the zones that contain the specified zone name
      --no-headers          Don't print table headers when table output is used
      --offset int32        Pagination offset
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --state string        Filter used to fetch all zones in a particular state (AVAILABLE, FAILED, PROVISIONING, DESTROYING)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dns z list
```

