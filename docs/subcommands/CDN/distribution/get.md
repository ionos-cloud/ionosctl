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
  -u, --api-url string           Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cdn' and env var 'IONOS_API_URL' (default "https://cdn.%s.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -i, --distribution-id string   The ID of the distribution you want to retrieve (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string          Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers               Don't print table headers when table output is used
      --offset int               pagination offset: Number of items to skip before starting to collect the results
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl cdn ds get --distribution-id ID
```

