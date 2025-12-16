---
description: "Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-dns-record"
---

# DnsRecordCreate

## Usage

```text
ionosctl dns record create [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `create` command:

```text
[c post]
```

## Description

Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-dns-record

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --content string    The content (Record Data) for your chosen record type. For example, if --type A, --content should be an IPv4 IP. (required)
      --enabled           When true - the record is visible for lookup (default true)
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name \*           The name of the DNS record.  Provide a wildcard i.e. \* to match requests for non-existent names under your DNS Zone name. Note that some terminals require '*' to be escaped, e.g. '\*' (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --priority int32    Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types.
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --ttl int32         Time to live. The amount of time the record can be cached by a resolver or server before it needs to be refreshed from the authoritative DNS server (default 3600)
  -t, --type string       Type of DNS Record. Can be one of: A, AAAA, CNAME, ALIAS, MX, NS, SRV, TXT, CAA, SSHFP, TLSA, SMIMEA, DS, HTTPS, SVCB, OPENPGPKEY, CERT, URI, RP, LOC (required) (default "AAAA")
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string       The ID or name of the DNS zone (required)
```

## Examples

```text
ionosctl dns r create --zone foo-bar.com --type A --content 1.2.3.4 --name \*
```

