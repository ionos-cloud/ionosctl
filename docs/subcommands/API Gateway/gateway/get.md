---
description: "Retrieve a gateway"
---

# ApigatewayGatewayGet

## Usage

```text
ionosctl apigateway gateway get [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a gateway

## Options

```text
  -u, --api-url string      Override default host URL (default "https://apigateway.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the gateway
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway gateway get --gateway-id GATEWAYID
```

