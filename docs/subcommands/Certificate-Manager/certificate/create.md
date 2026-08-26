---
description: "Upload a new certificate"
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

Upload a TLS/SSL certificate so IONOS products (ALB, CDN) can serve HTTPS under your domain.

You must supply three PEM values, each either inline or from a file (the inline and *-path variants are mutually exclusive, and exactly one of each pair is required):
  * the certificate body (--certificate / --certificate-path)
  * the certificate chain of intermediate CAs (--certificate-chain / --certificate-chain-path)
  * the private key that pairs with the certificate (--private-key / --private-key-path)

The private key is write-only; it is never returned by get/list. This command does not renew certificates - upload a fresh one before expiry. For automatic issuance and renewal, use `autocertificate create` instead.

## Options

```text
  -u, --api-url string                  Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --certificate string              The certificate body in PEM format, inline. Provide this or --certificate-path, not both
      --certificate-chain string        The chain of intermediate CA certificates in PEM format, inline. Provide this or --certificate-chain-path, not both
      --certificate-chain-path string   Path to a file containing the intermediate CA certificate chain in PEM format. Provide this or --certificate-chain, not both
  -n, --certificate-name string         A friendly name for the certificate, used to identify it in listings and when other products reference it (required)
      --certificate-path string         Path to a file containing the certificate body in PEM format. Provide this or --certificate, not both
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [CertId DisplayName Expired NotAfter NotBefore SerialNumber SubjectAlternativeNames Certificate CertificateChain]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                       Level of detail for response objects (default 1)
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Maximum number of items to return per request (default 50)
  -l, --location string                 Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --private-key string              The private key paired with the certificate, in PEM format, inline. Write-only: never returned by get/list. Provide this or --private-key-path, not both
      --private-key-path string         Path to a file containing the private key in PEM format. Provide this or --private-key, not both
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Upload from PEM files (recommended)
ionosctl certmanager certificate create --certificate-name my-cert --certificate-path cert.pem --certificate-chain-path chain.pem --private-key-path key.pem

# Upload with inline PEM values read from files by the shell
ionosctl certmanager certificate create --certificate-name my-cert --certificate "$(cat cert.pem)" --certificate-chain "$(cat chain.pem)" --private-key "$(cat key.pem)"
```

