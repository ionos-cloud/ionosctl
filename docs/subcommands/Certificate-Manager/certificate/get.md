---
description: "Get Certificate by ID"
---

# CertmanagerCertificateGet

## Usage

```text
ionosctl certmanager certificate get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve a Certificate by ID.

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --certificate             Print the certificate
      --certificate-chain       Print the certificate chain
  -i, --certificate-id string   Provide the specified Certificate (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName Expired NotAfter NotBefore]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers              Don't print table headers when table output is used
      --offset int              Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

