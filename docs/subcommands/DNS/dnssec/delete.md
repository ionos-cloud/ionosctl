---
description: "Removes ALL associated DNSKEY records for your DNS zone and disables DNSSEC keys."
---

# DnsDnssecDelete

## Usage

```text
ionosctl dns dnssec delete [flags]
```

## Aliases

For `dnssec` command:

```text
[sec dnskey key keys]
```

For `delete` command:

```text
[del rm remove]
```

## Description

Removes ALL associated DNSKEY records for your DNS zone and disables DNSSEC keys.

## Options

```text
  -u, --api-url string   Override default host url (default "dns.de-fra.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Id KeyTag DigestAlgorithmMnemonic Digest Validity Flags PubKey ComposedKeyData Algorithm KskBits ZskBits NsecMode Nsec3Iterations Nsec3SaltBits] (default [Id,KeyTag,DigestAlgorithmMnemonic,Digest,Validity])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
  -z, --zone string      The name or ID of the DNS zone
```

## Examples

```text
ionosctl dns keys delete --zone ZONE
```

