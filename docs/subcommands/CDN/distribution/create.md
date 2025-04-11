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
  -u, --api-url string          Override default host URL (default "https://cdn.de-fra.ionos.com")
      --certificate-id string   The ID of the certificate
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Domain CertificateId State]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --domain string           The domain of the distribution
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/fra
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --routing-rules string    The routing rules of the distribution. JSON string or file path of routing rules
      --routing-rules-example   Print an example of routing rules
  -t, --timeout duration        Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait                    Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl cdn ds create --domain foo-bar.com --certificate-id id --routing-rules rules.json
```

