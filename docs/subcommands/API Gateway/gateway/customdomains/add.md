---
description: ""
---

# ApigatewayGatewayCustomdomainsAdd

## Usage

```text
ionosctl apigateway gateway customdomains add [flags]
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
  -u, --api-url string          Override default host URL (default "https://apigateway.de-txl.ionos.com")
      --certificate-id string   The ID of the certificate to use for the distribution. (required)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [Id Name Logs Metrics Enable DomainName CertificateId HttpMethods HttpCodes Override PublicEndpoint Status]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -g, --gateway-id string       The ID of the gateway (required)
  -h, --help                    Print usage
  -l, --location string         Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit
      --name string             The domain name of the distribution. Field is validated as FQDN according to RFC1123. (required)
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway gateway customdomains add --gateway-id ID --name NAME --certificate-id CERTIFICATEID
```

