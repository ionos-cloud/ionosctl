---
description: "Update Certificate name"
---

# CertmanagerCertificateUpdate

## Usage

```text
ionosctl certmanager certificate update [flags]
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

For `update` command:

```text
[u]
```

## Description

Use this change a certificate's name.

## Options

```text
  -u, --api-url string            Override default host URL (default "https://dns.de-fra.ionos.com")
  -i, --certificate-id string     Provide certificate ID (required)
  -n, --certificate-name string   Provide new certificate name (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [CertId DisplayName]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -l, --location string           Location of the resource to operate on. Can be one of: de/fra
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -v, --verbose                   Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager update --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

