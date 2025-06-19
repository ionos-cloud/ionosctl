---
description: "Retrieve gateways"
---

# ApigatewayGatewayList

## Usage

```text
ionosctl apigateway gateway list [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve gateways

## Options

```text
  -u, --api-url string      Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://apigateway.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --max-results int32   Pagination limit
      --no-headers          Don't print table headers when table output is used
      --offset int32        Pagination offset
      --order-by string     The field to order the results by. If not provided, the results will be ordered by the default field.
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway gateway list
```

