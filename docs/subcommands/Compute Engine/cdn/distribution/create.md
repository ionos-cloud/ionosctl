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
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --certificate-id string   The ID of the certificate
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Domain CertificateId State]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --domain string           The domain of the distribution
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
      --routing-rules string    The routing rules of the distribution. JSON string or file path of routing rules
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionosctl cdn ds create --domain foo-bar.com --certificate-id id --routing-rules rules.json
```

