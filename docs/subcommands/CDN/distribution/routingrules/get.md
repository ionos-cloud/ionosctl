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
  -u, --api-url string           Override default host URL (default "https://cdn.de-fra.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --distribution-id string   The ID of the distribution (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
  -l, --location string          Location of the resource to operate on. Can be one of: de/fra
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -t, --timeout duration         Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                  Print step-by-step process when running command
  -w, --wait                     Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl cdn ds rr get --distribution-id ID
```

