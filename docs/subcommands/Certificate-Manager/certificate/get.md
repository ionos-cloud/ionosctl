---
description: "Get a certificate by ID"
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

Retrieve a single uploaded certificate by ID, including metadata such as its serial number, validity dates (notBefore/notAfter), and whether it has expired.

By default the metadata is printed as a table. To print the PEM material instead, pass --certificate (the certificate body) or --certificate-chain (the chain); these two are mutually exclusive. The private key is write-only and is never returned.

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --certificate             Instead of the metadata table, print only the certificate body PEM. Mutually exclusive with --certificate-chain
      --certificate-chain       Instead of the metadata table, print only the certificate chain PEM. Mutually exclusive with --certificate
  -i, --certificate-id string   The ID (UUID) of the certificate to retrieve (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName Expired NotAfter NotBefore SerialNumber SubjectAlternativeNames Certificate CertificateChain]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Show certificate metadata
ionosctl certmanager certificate get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422

# Print the certificate body PEM (e.g. to pipe into a file)
ionosctl certmanager certificate get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422 --certificate
```

