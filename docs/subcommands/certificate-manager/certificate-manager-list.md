---
description: List all Certificates
---

# CertificateManagerList

## Usage

```text
ionosctl certificate-manager list [flags]
```

## Aliases

For `list` command:

```text
[l]
```

## Description

Use this command to retrieve all Certificates.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Get response without headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager list
```

