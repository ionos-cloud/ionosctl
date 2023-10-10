---
description: "Get Certificate Manager API Version"
---

# CertificateManagerApiVersion

## Usage

```text
ionosctl certificate-manager api-version [flags]
```

## Aliases

For `api-version` command:

```text
[api info]
```

## Description

Use this command to retrieve API Version.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Name Href Version]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -n, --no-headers       Response delete all certificates
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager api-version
```

