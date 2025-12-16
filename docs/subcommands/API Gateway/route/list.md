---
description: "Retrieve routes"
---

# ApigatewayRouteList

## Usage

```text
ionosctl apigateway route list [flags]
```

## Aliases

For `route` command:

```text
[r]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve routes

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'apigateway' and env var 'IONOS_API_URL' (default "https://apigateway.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Type Paths Methods Host Port Weight Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the gateway (required)
  -h, --help                Print usage
      --limit int32         Pagination limit
  -l, --location string     Location of the resource to operate on. Can be one of: de/txl, gb/lhr, fr/par, es/vit (default "de/txl")
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     The field to order the results by. If not provided, the results will be ordered by the default field.
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl apigateway route list --gateway-id ID
```

