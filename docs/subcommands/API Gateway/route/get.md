---
description: "Retrieve a route"
---

# ApigatewayRouteGet

## Usage

```text
ionosctl apigateway route get [flags]
```

## Aliases

For `route` command:

```text
[r]
```

For `get` command:

```text
[g]
```

## Description

Retrieve a route

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --route-id string     The ID of the route (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl apigateway route get --gateway-id GATEWAYID --route-id ROUTEID
```

