---
description: "Enable DNSSEC keys and create associated DNSKEY records for your DNS zone"
---

# DnsDnssecCreate

## Usage

```text
ionosctl dns dnssec create [flags]
```

## Aliases

For `dnssec` command:

```text
[sec dnskey key keys]
```

For `create` command:

```text
[c post]
```

## Description

Enable DNSSEC keys and create associated DNSKEY records for your DNS zone

## Options

```text
      --algorithm string       Algorithm used to generate signing keys (both Key Signing Keys and Zone Signing Keys) (default "RSASHA256")
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Id KeyTag DigestAlgorithmMnemonic Digest Validity Flags PubKey ComposedKeyData Algorithm KskBits ZskBits NsecMode Nsec3Iterations Nsec3SaltBits] (default [Id,KeyTag,DigestAlgorithmMnemonic,Digest,Validity])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ksk-bits int           Key signing key length in bits. kskBits >= zskBits: [1024/2048/4096] (default 1024)
      --limit int              pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers             Don't print table headers when table output is used
      --nsec-mode string       NSEC mode.. Can be one of: NSEC, NSEC3 (default "NSEC")
      --nsec3-iterations int   Number of iterations for NSEC3. [0..50]
      --nsec3-salt-bits int    Salt length in bits for NSEC3. [64..128], multiples of 8 (default 64)
      --offset int             pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --validity int           Signature validity in days [90..365] (default 90)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -z, --zone string            The name or ID of the DNS zone (required)
      --zsk-bits int           Zone signing key length in bits. zskBits <= kskBits: [1024/2048/4096] (default 1024)
```

## Examples

```text
ionosctl dns keys create --zone ZONE
```

