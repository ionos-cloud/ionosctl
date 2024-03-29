---
description: "Add a new Certificate"
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
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --certificate string              Specify the certificate itself (required either this or --certificate-path)
      --certificate-chain string        Specify the certificate chain (required either this or --certificate-chain-path)
      --certificate-chain-path string   Specify the certificate chain from a file (required either this or --certificate-chain)
  -n, --certificate-name string         Specify name of the certificate (required)
      --certificate-path string         Specify the certificate itself from a file (required either this or --certificate)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [CertId DisplayName]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --no-headers                      Don't print table headers when table output is used
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --private-key string              Specify the private key (required either this or --private-key-path)
      --private-key-path string         Specify the private key from a file (required either this or --private-key)
  -q, --quiet                           Quiet output
  -v, --verbose                         Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager add --certificate-name my-cert --certificate-path /path/to/cert --certificate-chain-path /path/to/cert-chain --private-key-path /path/to/private-key
ionosctl certificate-manager add --certificate-name my-cert --certificate <certificate> --certificate-chain <certificate chain> --private-key <private key>
```

