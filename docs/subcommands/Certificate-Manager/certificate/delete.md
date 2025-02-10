---
description: "Delete Certificate by ID or all Certificates"
---

# CertmanagerCertificateDelete

## Usage

```text
ionosctl certmanager certificate delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a Certificate by ID.

## Options

```text
  -a, --all                     Response delete all certificates
  -u, --api-url string          Override default host URL (default "https://dns.de-fra.ionos.com")
  -i, --certificate-id string   Response delete a single certificate (required)
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
ionosctl certificate-manager delete --certificate-id CERTIFICATE_ID 
ionosctl certificate-manager delete --all
```

