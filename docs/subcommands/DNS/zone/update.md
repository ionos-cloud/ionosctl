---
description: "Partially modify a zone's properties. This command uses a combination of GET and PUT to simulate a PATCH operation"
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
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Description NameServers Enabled State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string   The new description of the DNS zone
      --enabled              Activate or deactivate the DNS zone (default true)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string          The new name of the DNS zone, e.g. foo.com
      --no-headers           Don't print table headers when table output is used
      --offset int           pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string          The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns z update --zone ZONE --name newname.com
```

