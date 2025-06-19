---
description: "Delete Certificate by ID or all Certificates"
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
  -i, --certificate-id string   Response delete a single certificate (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [CertId DisplayName]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionsoctl certificate-manager delete --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

