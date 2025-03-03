---
description: "Get Certificate by ID"
---

# CertificateManagerGet

## Usage

```text
ionosctl certificate-manager get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve a Certificate by ID.

## Options

```text
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --certificate             Print the certificate
      --certificate-chain       Print the certificate chain
  -i, --certificate-id string   Response get a single certificate (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose count           Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

