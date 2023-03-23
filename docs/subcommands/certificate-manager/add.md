---
description: Add a new Certificate
---

# CertificateManagerAdd

## Usage

```text
ionosctl certificate-manager add [flags]
```

## Aliases

For `add` command:

```text
[a]
```

## Description

Use this command to add a Certificate.

## Options

```text
      --certificate string              Specify the certificate itself (required either this or --certificate-path)
      --certificate-chain string        Specify the certificate chain (required either this or --certificate-chain-path)
      --certificate-chain-path string   Specify the certificate chain from a file (required either this or --certificate-chain)
  -n, --certificate-name string         Specify name of the certificate (required)
      --certificate-path string         Specify the certificate itself from a file (required either this or --certificate)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [CertId DisplayName]
      --private-key string              Specify the private key (required either this or --private-key-path)
      --private-key-path string         Specify the private key from a file (required either this or --private-key)
```

## Examples

```text
ionosctl certificate-manager add --certificate-name my-cert --certificate-path /path/to/cert --certificate-chain-path /path/to/cert-chain --private-key-path /path/to/private-key
ionosctl certificate-manager add --certificate-name my-cert --certificate <certificate> --certificate-chain <certificate chain> --private-key <private key>
```

