---
description: "Delete an AutoCertificate"
---

# CertmanagerAutocertificateDelete

## Usage

```text
ionosctl certmanager autocertificate delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete an AutoCertificate

## Options

```text
  -a, --all                         Delete all AutoCertificates
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
  -i, --autocertificate-id string   Provide the specified AutoCertificate (required)
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Id Provider CommonName KeyAlgorithm Name AlternativeNames State]
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
      --limit int                   Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string             Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string               Desired output format [text|json|api-json] (default "text")
  -q, --quiet                       Quiet output
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certmanager autocertificate delete --autocertificate-id ID
```

