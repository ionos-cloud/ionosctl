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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --max-results int32   Pagination limit
  -n, --name string         Filter used to fetch only the zones that contain the specified zone name
      --no-headers          Don't print table headers when table output is used
      --offset int32        Pagination offset
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --state string        Filter used to fetch all zones in a particular state (AVAILABLE, FAILED, PROVISIONING, DESTROYING)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dns z list
```

