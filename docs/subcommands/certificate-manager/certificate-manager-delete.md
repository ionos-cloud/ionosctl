---
description: Delete Certificate by ID or all Certificates
---

# CertificateManagerDelete

## Usage

```text
ionosctl certificate-manager delete [flags]
```

## Aliases

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a Certificate by ID.

## Options

```text
  -a, --all                     Response delete all certificates
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
  -i, --certificate-id string   Response delete a single certificate (required)
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -o, --output string           Desired output format [text|json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionsoctl certificate-manager delete --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

