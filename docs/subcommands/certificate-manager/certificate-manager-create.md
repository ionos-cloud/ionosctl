---
description: Add a new Certificate
---

# CertificateManagerCreate

## Usage

```text
ionosctl certificate-manager create [flags]
```

## Aliases

For `create` command:

```text
[c]
```

## Description

Use this command to add a Certificate.

## Options

```text
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --certificate string              Specify the certificate itself (required either this or --certificate-path)
      --certificate-chain string        Specify the certificate chain (required either this or --certificate-chain-path)
      --certificate-chain-path string   Specify the certificate chain from a file (required either this or --certificate-chain)
  -n, --certificate-name string         Specify name of the certificate (required)
      --certificate-path string         Specify the certificate itself from a file (required either this or --certificate)
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --no-headers                      Get response without headers
  -o, --output string                   Desired output format [text|json] (default "text")
      --private-key string              Specify the private key (required either this or --private-key-path)
      --private-key-path string         Specify the private key from a file (required either this or --private-key)
  -q, --quiet                           Quiet output
  -v, --verbose                         Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager create
```

