---
description: "Partially modify a distribution's properties. This command uses a combination of GET and PUT to simulate a PATCH operation"
---

# CdnDistributionUpdate

## Usage

```text
ionosctl cdn distribution update [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `update` command:

```text
[u]
```

## Description

Partially modify a distribution's properties. This command uses a combination of GET and PUT to simulate a PATCH operation

## Options

```text
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --certificate-id string    The ID of the certificate
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [Id Domain CertificateId State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --distribution-id string   The ID of the CDN distribution (required)
      --domain string            The domain of the distribution
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
      --routing-rules string     The routing rules of the distribution. JSON string or file path of routing rules
      --routing-rules-example    Print an example of routing rules
  -v, --verbose                  Print step-by-step process when running command
```

## Examples

```text
ionosctl cdn ds update --distribution-id
```

