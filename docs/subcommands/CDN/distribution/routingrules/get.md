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
  -u, --api-url string           Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cdn' and env var 'IONOS_API_URL' (default "https://cdn.%s.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                Level of detail for response objects (default 1)
  -i, --distribution-id string   The ID of the distribution (required)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
  -l, --location string          Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl cdn ds rr get --distribution-id ID
```

