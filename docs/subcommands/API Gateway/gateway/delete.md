---
description: "Delete a gateway"
---

# ApigatewayGatewayDelete

## Usage

```text
ionosctl apigateway gateway delete [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `delete` command:

```text
[del d]
```

## Description

Delete a gateway

## Options

```text
  -a, --all                 Delete all gateways. Required or -g
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the gateway. Required or -a
  -h, --help                Print usage
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl apigateway gateway delete --gateway-id ID
```

