---
description: "Delete a record"
---

# DnsReverseRecordDelete

## Usage

```text
ionosctl dns reverse-record delete [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `delete` command:

```text
[d del rm]
```

## Description

Delete a record

## Options

```text
  -a, --all               Delete all records if set (required)
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
      --record string     The record ID or IP which you want to delete (required)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl dns rr delete -af
ionosctl dns rr delete --record RECORD_IP
ionosctl dns rr delete --record RECORD_ID
```

