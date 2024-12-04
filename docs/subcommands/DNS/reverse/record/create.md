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
  -u, --api-url string       Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name IP Description] (default [Id,Name,IP,Description])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --description string   Description stored along with the reverse DNS record to describe its usage
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --ip string            [IPv4/IPv6] Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract (required)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra
  -n, --name string          The name of the DNS Reverse Record. (required)
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dns rr create --name mail.example.com --ip 5.6.7.8
```

