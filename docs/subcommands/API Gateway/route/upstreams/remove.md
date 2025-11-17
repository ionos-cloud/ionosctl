---
description: "Upstreams consist of schme, loadbalancer, host, port and weight"
---

# ApigatewayRouteUpstreamsRemove

## Usage

```text
ionosctl apigateway route upstreams remove [flags]
```

## Aliases

For `route` command:

```text
[r]
```

For `upstreams` command:

```text
[streams]
```

For `remove` command:

```text
[r]
```

## Description

Upstreams consist of schme, loadbalancer, host, port and weight

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -i, --gateway-id string    The ID of the gateway (required)
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
      --route-id string      The ID of the route. Required or -a (required)
      --upstream-id string   The ID of the upstream. Required or -a (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl apigateway route upstreams remove --gateway-id ID --route-id ID_ROUTE --upstream-id UPSTREAMID
```

