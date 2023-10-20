---
description: "Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation"
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

Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation

## Options

```text
  -u, --api-url string   Override default host url (default "dns.de-fra.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id Name Content Type Enabled FQDN State ZoneId ZoneName]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --content string   The content (Record Data) for your chosen record type. For example, if --type A, --content should be an IPv4 IP. (required)
      --enabled          When true - the record is visible for lookup (default true)
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --name \*          The name of the DNS record.  Provide a wildcard i.e. \* to match requests for non-existent names under your DNS Zone name. Note that some terminals require '*' to be escaped, e.g. '\*' (required)
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
      --priority int32   Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types.
  -q, --quiet            Quiet output
  -r, --record string    The ID or name of the DNS record (required)
      --ttl int32        Time to live. The amount of time the record can be cached by a resolver or server before it needs to be refreshed from the authoritative DNS server (default 3600)
  -t, --type string      Type of DNS Record. Can be one of: A, AAAA, CNAME, ALIAS, MX, NS, SRV, TXT, CAA, SSHFP, TLSA, SMIMEA, DS, HTTPS, SVCB, OPENPGPKEY, CERT, URI, RP, LOC (required) (default "AAAA")
  -v, --verbose          Print step-by-step process when running command
  -z, --zone string      The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns z update --zone ZONE_ID --record RECORD
```

