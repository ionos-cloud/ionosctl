---
description: Partially modify a zone's properties. This command uses a combination of GET and PUT to simulate a PATCH operation
---

# DnsZoneUpdate

## Usage

```text
ionosctl dns zone update [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `update` command:

```text
[u]
```

## Description

Partially modify a zone's properties. This command uses a combination of GET and PUT to simulate a PATCH operation

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --description string   The description of the DNS zone
      --enabled              Activate or deactivate the DNS zone (default true)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -n, --name string          The name of the DNS zone, e.g. foo.com
      --no-headers           When using text output, don't print headers
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
  -z, --zone string          The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns zone update --zone ZONE --name newname.com
```
