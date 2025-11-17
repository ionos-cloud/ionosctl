---
description: "Create an AutoCertificate. Requires an enabled DNS Zone with the same name as the --common-name."
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

Create an AutoCertificate. Requires an enabled DNS Zone with the same name as the --common-name.

## Options

```text
  -u, --api-url string                      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Id Provider CommonName KeyAlgorithm Name AlternativeNames State]
      --common-name string                  The common name (DNS) of the certificate to issue
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                           Level of detail for response objects (default 1)
      --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --key-algorithm string                The key algorithm used to generate the certificate. (required)
      --limit int                           Maximum number of items to return per request (default 50)
  -l, --location string                     Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string                         The name of the AutoCertificate
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -i, --provider-id string                  The certificate provider used to issue the AutoCertificate (required)
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
      --subject-alternative-names strings   Optional additional names to be added to the issued certificate
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certmanager autocertificate create --name NAME --provider-id PROVIDER --common-name COMMONNAME --key-algorithm rsa2048
```

