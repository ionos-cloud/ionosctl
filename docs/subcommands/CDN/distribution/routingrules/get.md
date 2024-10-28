---
description: "Retrieve a distribution routing rules"
---

# CdnDistributionRoutingrulesGet

## Usage

```text
ionosctl cdn distribution routingrules get [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `routingrules` command:

```text
[rr]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a distribution routing rules

## Options

```text
  -u, --api-url string           Override default host url (default "https://cdn.de-fra.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --distribution-id string   The ID of the distribution you want to retrieve (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose                  Print step-by-step process when running command
```

## Examples

```text
ionosctl cdn ds rr get --distribution-id ID
```

