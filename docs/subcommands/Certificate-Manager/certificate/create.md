---
description: "Add a new Certificate"
---

# CertmanagerCertificateCreate

## Usage

```text
ionosctl certmanager certificate create [flags]
```

## Aliases

For `certmanager` command:

```text
[cert certs certificate-manager certificates certificate]
```

For `certificate` command:

```text
[cert certificates certs]
```

For `create` command:

```text
[add a c post]
```

## Description

Use this command to add a Certificate.

## Options

```text
  -u, --api-url string                  Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --certificate string              Specify the certificate itself (required either this or --certificate-path)
      --certificate-chain string        Specify the certificate chain (required either this or --certificate-chain-path)
      --certificate-chain-path string   Specify the certificate chain from a file (required either this or --certificate-chain)
  -n, --certificate-name string         Specify name of the certificate (required)
      --certificate-path string         Specify the certificate itself from a file (required either this or --certificate)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [CertId DisplayName Expired NotAfter NotBefore]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string                 Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers                      Don't print table headers when table output is used
      --offset int                      pagination offset: Number of items to skip before starting to collect the results
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --private-key string              Specify the private key (required either this or --private-key-path)
      --private-key-path string         Specify the private key from a file (required either this or --private-key)
  -q, --quiet                           Quiet output
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certmanager certificate create --certificate-name CERTIFICATE_NAME --certificate-chain-path CERTIFICATE_CHAIN_PATH --certificate-path CERTIFICATE_PATH --private-key-path PRIVATE_KEY_PATH 
```

