---
description: "Retrieve your zone's DNSSEC keys"
---

# DnsDnssecList

## Usage

```text
ionosctl dns dnssec list [flags]
```

## Aliases

For `dnssec` command:

```text
[sec dnskey key keys]
```

For `list` command:

```text
[l ls get g]
```

## Description

Retrieve your zone's DNSSEC keys

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id KeyTag DigestAlgorithmMnemonic Digest Validity Flags PubKey ComposedKeyData Algorithm KskBits ZskBits NsecMode Nsec3Iterations Nsec3SaltBits] (default [Id,KeyTag,DigestAlgorithmMnemonic,Digest,Validity])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
      --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string       The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns keys list --zone ZONE
ionosctl dns keys list --zone ZONE --cols ComposedKeyData --no-headers
ionosctl dns keys list --zone ZONE --cols PubKey --no-headers
```

