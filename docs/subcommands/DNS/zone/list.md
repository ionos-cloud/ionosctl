---
description: Retrieve zones
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --max-results int32   Pagination limit
      --name string         Filter used to fetch only the zones that contain the specified zone name
      --offset int32        Pagination offset
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --state string        Filter used to fetch all zones in a particular state (PROVISIONING, DEPROVISIONING, CREATED, FAILED)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dns zone list
```

