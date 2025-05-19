---
description: "Create a zone"
---

# DnsZoneCreate

## Usage

```text
ionosctl dns zone create [flags]
```

## Aliases

For `zone` command:

```text
[z zones]
```

For `create` command:

```text
[post c]
```

## Description

Create a zone

## Options

```text
  -u, --api-url string       Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string   The description of the DNS zone
      --enabled              Activate or deactivate the DNS zone (default true)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra
  -n, --name string          The name of the DNS zone, e.g. foo.com
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dns z create --name name.com
```

