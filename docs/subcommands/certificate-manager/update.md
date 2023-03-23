---
description: Update Certificate name
---

# CertificateManagerUpdate

## Usage

```text
ionosctl certificate-manager update [flags]
```

## Aliases

For `update` command:

```text
[u]
```

## Description

Use this change a certificate's name.

## Options

```text
  -i, --certificate-id string     Provide certificate ID (required)
  -n, --certificate-name string   Provide new certificate name (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [CertId DisplayName]
```

## Examples

```text
ionosctl certificate-manager update --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

