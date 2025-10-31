---
description: "Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-and-manage-reverse-dns"
---

# DnsReverseRecordCreate

## Usage

```text
ionosctl dns reverse-record create [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `create` command:

```text
[c post]
```

## Description

Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-and-manage-reverse-dns

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string   Description stored along with the reverse DNS record to describe its usage
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --ip string            [IPv4/IPv6] Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract (required)
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string          The name of the DNS Reverse Record. (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dns rr create --name mail.example.com --ip 5.6.7.8
```

