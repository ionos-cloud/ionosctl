---
description: "Issue an auto-renewing certificate via a provider"
---

# CertmanagerAutocertificateCreate

## Usage

```text
ionosctl certmanager autocertificate create [flags]
```

## Aliases

For `certmanager` command:

```text
[cert certs certificate-manager certificates certificate]
```

For `autocertificate` command:

```text
[autocert autocerts auto autocertificates]
```

For `create` command:

```text
[post c]
```

## Description

Request a certificate that IONOS issues and auto-renews through an ACME provider.

IONOS validates domain ownership with ACME DNS-01 challenges, so the --common-name (and every --subject-alternative-names entry) must belong to a DNS zone hosted in IONOS CLOUD DNS that you manage; IONOS creates the required TXT records automatically. If the matching zone does not exist, issuance fails.

Required: --name, --provider-id (an existing provider, see 'certmanager provider create'), --common-name, and --key-algorithm. Once issued, the certificate renews automatically ~30 days before expiry.

## Options

```text
  -u, --api-url string                      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Id Provider CommonName KeyAlgorithm Name AlternativeNames State]
      --common-name string                  The primary domain (DNS name) the certificate is issued for, e.g. www.example.com. Must belong to an IONOS CLOUD DNS zone you manage. Required
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --key-algorithm string                The key algorithm for the generated private key. One of: rsa2048, rsa3072, rsa4096. Required (required)
      --limit int                           Maximum number of items to return per request (default 50)
  -l, --location string                     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string                         A friendly name to identify the auto-certificate for management purposes. Required
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -i, --provider-id string                  The ID (UUID) of the ACME provider that will issue and renew this certificate (required)
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
      --subject-alternative-names strings   Additional domains (SANs) to cover with the same certificate, comma-separated. Each must also belong to an IONOS CLOUD DNS zone you manage. Optional
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Issue a single-domain RSA-2048 certificate
ionosctl certmanager autocertificate create --name web-cert --provider-id PROVIDER_ID --common-name www.example.com --key-algorithm rsa2048

# Issue a certificate covering additional (SAN) domains with a stronger key
ionosctl certmanager autocertificate create --name web-cert --provider-id PROVIDER_ID --common-name www.example.com --key-algorithm rsa4096 --subject-alternative-names app.example.com,api.example.com
```

