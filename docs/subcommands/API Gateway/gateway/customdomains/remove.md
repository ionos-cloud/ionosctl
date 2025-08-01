---
description: ""
---

# ApigatewayGatewayCustomdomainsRemove

## Usage

```text
ionosctl apigateway gateway customdomains remove [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `customdomains` command:

```text
[custom-domains]
```

## Description

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --custom-domains-id string   The ID of the custom domain (required)
  -f, --force                      Force command to execute without user input
  -g, --gateway-id string          The ID of the gateway (required)
  -h, --help                       Print usage
  -l, --location string            Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers                 Don't print table headers when table output is used
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                      Quiet output
  -v, --verbose                    Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway gateway customdomains remove --gateway-id ID --custom-domains-id ID
```

