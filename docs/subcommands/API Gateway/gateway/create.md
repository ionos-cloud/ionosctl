---
description: "Create an apigateway"
---

# ApigatewayGatewayCreate

## Usage

```text
ionosctl apigateway gateway create [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `create` command:

```text
[post c]
```

## Description

Create an apigateway

## Options

```text
  -u, --api-url string                         Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://apigateway.de-txl.ionos.com")
      --cols strings                           Set of columns to be printed on output 
                                               Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string                          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --custom-domains-certificate-id string   The ID of the certificate to use for the distribution.
      --custom-domains-name string             The domain name of the distribution. Field is validated as FQDN
  -f, --force                                  Force command to execute without user input
  -h, --help                                   Print usage
  -l, --location string                        Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --logs                                   The logs parameter of the ApiGateway gateway
      --metrics                                Activate or deactivate the ApiGateway gateway metrics parameter
  -n, --name string                            The name of the ApiGateway gateway
      --no-headers                             Don't print table headers when table output is used
  -o, --output string                          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                                  Quiet output
  -v, --verbose                                Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway gateway create --name name
```

