---
description: "List secondary zones"
---

# DnsSecondaryZoneList

## Usage

```text
ionosctl dns secondary-zone list [flags]
```

## Aliases

For `secondary-zone` command:

```text
[secondary-zones sz]
```

## Description

List all secondary zones. Default limit is the first 100 items. Use pagination query parameters for listing more items (up to 1000).

## Options

```text
  -u, --api-url string      Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Description PrimaryIPs State] (default [Id,Name,Description,PrimaryIPs,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra
      --max-results int32   Pagination limit
  -n, --name string         Filter used to fetch only the zones that contain the specified zone name
      --no-headers          Don't print table headers when table output is used
      --offset int32        Pagination offset
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --state string        Filter used to fetch all zones in a particular state (AVAILABLE, FAILED, PROVISIONING, DESTROYING)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dns secondary-zone list
```

