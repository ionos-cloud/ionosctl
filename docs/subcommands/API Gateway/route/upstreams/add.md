---
description: "Upstreams consist of schme, loadbalancer, host, port and weight"
---

# ApigatewayRouteUpstreamsAdd

## Usage

```text
ionosctl apigateway route upstreams add [flags]
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

For `add` command:

```text
[a]
```

## Description

Upstreams consist of schme, loadbalancer, host, port and weight

## Options

```text
  -u, --api-url string        Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                 Force command to execute without user input
  -i, --gateway-id string     The ID of the gateway (required)
  -h, --help                  Print usage
      --host string           The host of the upstream. Field is validated as hostname according to RFC1123. (required)
      --limit int             Pagination limit: Maximum number of items to return per request (default 50)
      --loadbalancer string   The load balancer algorithm. (default "roundrobin")
  -l, --location string       Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers            Don't print table headers when table output is used
      --offset int            Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --port int32            The port of the upstream. (default 80)
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
      --route-id string       The ID of the route. Required or -a
  -s, --scheme string         The target URL of the upstream.. Can be one of: http, https, grpc, grpcs (default "http")
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
      --weight int32          Weight with which to split traffic to the upstream. (default 100)
```

## Examples

```text
ionosctl apigateway route upstreams add --gateway-id ID --route-id ID_ROUTE --host HOST
```

