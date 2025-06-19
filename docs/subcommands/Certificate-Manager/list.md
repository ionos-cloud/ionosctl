---
description: "List all Certificates"
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
      --cols strings    Set of columns to be printed on output 
                        Available columns: [CertId DisplayName]
  -c, --config string   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force           Force command to execute without user input
  -h, --help            Print usage
      --no-headers      Don't print table headers when table output is used
  -o, --output string   Desired output format [text|json|api-json] (default "text")
  -q, --quiet           Quiet output
  -v, --verbose         Print step-by-step process when running command
```

## Examples

```text
ionosctl certificate-manager list
```

