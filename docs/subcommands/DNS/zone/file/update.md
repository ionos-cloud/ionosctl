---
description: "Update a zone file"
---

# DnsZoneFileUpdate

## Usage

```text
ionosctl dns zone file update [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `file` command:

```text
[f]
```

## Description

Update a zone file

## Options

```text
  -u, --api-url string     Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -l, --location string    Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose            Print step-by-step process when running command
  -z, --zone string        The name or ID of the DNS zone (required)
      --zone-file string   Path to the zone file
```

## Examples

```text
ionosctl dns zone file update --zone ZONE_ID --zone-file FILE_PATH
```

