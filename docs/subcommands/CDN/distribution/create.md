---
description: "Create a CDN distribution. Wiki: https://docs.ionos.com/cloud/network-services/cdn/dcd-how-tos/create-cdn-distribution"
---

# CdnDistributionCreate

## Usage

```text
ionosctl cdn distribution create [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `create` command:

```text
[c post]
```

## Description

Create a CDN distribution. Wiki: https://docs.ionos.com/cloud/network-services/cdn/dcd-how-tos/create-cdn-distribution

## Options

```text
  -u, --api-url string          Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cdn' and env var 'IONOS_API_URL' (default "https://cdn.%s.ionos.com")
      --certificate-id string   The ID of the certificate
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Domain CertificateId State]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --domain string           The domain of the distribution
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers              Don't print table headers when table output is used
      --offset int              Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --routing-rules string    The routing rules of the distribution. JSON string or file path of routing rules
      --routing-rules-example   Print an example of routing rules
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl cdn ds create --domain foo-bar.com --certificate-id id --routing-rules rules.json
```

