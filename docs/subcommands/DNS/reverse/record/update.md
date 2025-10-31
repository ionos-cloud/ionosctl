---
description: "Update a record"
---

# DnsReverseRecordUpdate

## Usage

```text
ionosctl dns reverse-record update [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `update` command:

```text
[u up]
```

## Description

Update a record

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string   The new description of the record
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --ip string            The new IP
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string          The new record name
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
      --record string        The record ID or IP which you want to update (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dns rr update --record OLD_RECORD_IP --name mail.example.com --ip 5.6.7.8
```

