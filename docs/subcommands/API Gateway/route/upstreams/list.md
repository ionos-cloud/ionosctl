---
description: "Upstreams consist of schme, loadbalancer, host, port and weight"
---

# ApigatewayRouteUpstreamsList

## Usage

```text
ionosctl apigateway route upstreams list [flags]
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

For `list` command:

```text
[l]
```

## Description

Upstreams consist of schme, loadbalancer, host, port and weight

## Options

```text
  -u, --api-url string      Override default host URL (default "https://apigateway.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -g, --gateway-id string   The ID of the gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --route-id string     The ID of the route. Required or -a
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway route upstreams list --gateway-id ID --route-id ID_ROUTE
```

