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
  -u, --api-url string    Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id KeyTag DigestAlgorithmMnemonic Digest Validity Flags PubKey ComposedKeyData Algorithm KskBits ZskBits NsecMode Nsec3Iterations Nsec3SaltBits] (default [Id,KeyTag,DigestAlgorithmMnemonic,Digest,Validity])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for polling the request (default 60)
  -v, --verbose           Print step-by-step process when running command
  -w, --wait              Polls the request continuously until the operation is completed 
  -z, --zone string       The name or ID of the DNS zone (required)
```

## Examples

```text
ionosctl dns keys delete --zone ZONE
```

