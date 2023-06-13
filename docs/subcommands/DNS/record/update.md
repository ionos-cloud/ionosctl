---
description: Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation
---

# DnsRecordUpdate

## Usage

```text
ionosctl dns record update [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `update` command:

```text
[u]
```

## Description

Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation.
You must use either --zone-id and --record-id, or alternatively use filters: --name and/or --zone-id. Note that if choosing to use filters, the operation will fail if more than one record is found

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [Id Name Content Type Enabled FQDN State ZoneId ZoneName]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --content string     The content (Record Data) for your chosen record type. For example, if --type A, --content should be an IPv4 IP. (required)
      --enabled            When true - the record is visible for lookup (default true)
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -n, --name *             The name of the DNS record.  Provide a wildcard i.e. * to match requests for non-existent names under your DNS Zone name (required)
      --no-headers         When using text output, don't print headers
  -o, --output string      Desired output format [text|json] (default "text")
      --priority int32     Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types.
  -q, --quiet              Quiet output
  -i, --record-id string   The ID (UUID) of the DNS record (required)
      --ttl int32          Time to live. The amount of time the record can be cached by a resolver or server before it needs to be refreshed from the authoritative DNS server (default 3600)
  -t, --type string        Type of DNS Record. Can be one of: A, AAAA, CNAME, ALIAS, MX, NS, SRV, TXT, CAA, SSHFP, TLSA, SMIMEA, DS, HTTPS, SVCB, OPENPGPKEY, CERT, URI, RP, LOC (required) (default "AAAA")
  -v, --verbose            Print step-by-step process when running command
      --zone-id string     The ID (UUID) of the DNS zone (required)
```

## Examples

```text
ionosctl dns zone update --zone ZONE --record-id RECORD_ID
```

