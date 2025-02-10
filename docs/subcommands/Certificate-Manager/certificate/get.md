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
  -u, --api-url string          Override default host URL (default "https://dns.de-fra.ionos.com")
      --certificate             Print the certificate
      --certificate-chain       Print the certificate chain
  -i, --certificate-id string   Response get a single certificate (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

