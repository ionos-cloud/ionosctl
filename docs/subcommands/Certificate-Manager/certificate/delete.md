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
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
  -i, --certificate-id string   Provide the specified Certificate (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName Expired NotAfter NotBefore]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certificate-manager delete --certificate-id CERTIFICATE_ID 
ionosctl certificate-manager delete --all
```

