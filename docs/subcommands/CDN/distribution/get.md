---
description: "Retrieve a distribution"
---

# CdnDistributionGet

## Usage

```text
ionosctl cdn distribution get [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a distribution

## Options

```text
  -u, --api-url string           Override default host URL (default "https://cdn.de-fra.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --distribution-id string   The ID of the distribution you want to retrieve (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
  -l, --location string          Location of the resource to operate on. Can be one of: de/fra
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose                  Print step-by-step process when running command
```

## Examples

```text
ionosctl cdn ds get --distribution-id ID
```

