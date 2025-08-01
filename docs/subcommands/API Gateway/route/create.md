---
description: "Create a new route"
---

# ApigatewayRouteCreate

## Usage

```text
ionosctl apigateway route create [flags]
```

## Aliases

For `route` command:

```text
[r]
```

For `create` command:

```text
[c post]
```

## Description

Create a new route

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
      --loadbalancer string   The load balancer algorithm. (default "roundrobin")
  -l, --location string       Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
  -m, --methods strings       The HTTP methods that the route should match. (required)
  -n, --name string           The name of the route. (required)
      --no-headers            Don't print table headers when table output is used
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --paths string          The paths that the route should match. (required)
      --port int32            The port of the upstream. (default 80)
  -q, --quiet                 Quiet output
  -s, --scheme string         The target URL of the upstream.. Can be one of: http, https, grpc, grpcs (required) (default "http")
      --type string            Default: http. This field specifies the protocol used by the ingress to route traffic to the backend service. (default "http")
  -v, --verbose               Print step-by-step process when running command
      --websocket             To enable websocket support.
      --weight int32          Weight with which to split traffic to the upstream. (default 100)
```

## Examples

```text
ionosctl apigateway route create --gateway-id ID --name NAME --paths PATHS --methods METHODS --host HOST
```

