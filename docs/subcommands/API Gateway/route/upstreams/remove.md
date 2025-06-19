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
  -u, --api-url string       Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://apigateway.de-txl.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -i, --gateway-id string    The ID of the gateway (required)
  -h, --help                 Print usage
  -l, --location string      Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
      --route-id string      The ID of the route. Required or -a (required)
      --upstream-id string   The ID of the upstream. Required or -a (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway route upstreams remove --gateway-id ID --route-id ID_ROUTE --upstream-id UPSTREAMID
```

