---
description: "Partially modify a gateway's properties. This command uses a combination of GET and PUT to simulate a PATCH operation"
---

# ApigatewayGatewayUpdate

## Usage

```text
ionosctl apigateway gateway update [flags]
```

## Aliases

For `gateway` command:

```text
[a api]
```

For `update` command:

```text
[u]
```

## Description

Partially modify a gateway's properties. This command uses a combination of GET and PUT to simulate a PATCH operation

## Options

```text
  -u, --api-url string                         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings                           Set of columns to be printed on output 
                                               Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string                          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --custom-domains-certificate-id string   The ID of the certificate to use for the distribution.
      --custom-domains-name string             The domain name of the distribution. Field is validated as FQDN according to RFC1123.
  -D, --depth int                              Level of detail for response objects (default 1)
      --filters strings                        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                  Force command to execute without user input
  -g, --gateway-id string                      The ID of the gateway (required)
  -h, --help                                   Print usage
      --limit int                              Maximum number of items to return per request (default 50)
  -l, --location string                        Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --logs                                   This field enables or disables the collection and reporting of logs for observability of this instance.
      --metrics                                This field enables or disables the collection and reporting of metrics for observability of this instance.
  -n, --name string                            The new name of the ApiGateway (required)
      --no-headers                             Don't print table headers when table output is used
      --offset int                             Number of items to skip before starting to collect the results
      --order-by string                        Property to order the results by
  -o, --output string                          Desired output format [text|json|api-json] (default "text")
      --query string                           JMESPath query string to filter the output
  -q, --quiet                                  Quiet output
  -v, --verbose count                          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl apigateway gateway update --gateway-id GATEWAYID --logs true
```

