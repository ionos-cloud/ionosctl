---
description: "Remove a NAT Gateway Lan"
---

# NatgatewayLanRemove

## Usage

```text
ionosctl natgateway lan remove [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a specified NAT Gateway Lan from a NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id

## Options

```text
  -a, --all                    Remove all NAT Gateway Lans.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps] (default [NatGatewayLanId,GatewayIps])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --lan-id int             The unique LAN Id (required) (default 1)
      --limit int              Maximum number of items to return per request (default 50)
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for NAT Gateway Lan deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for NAT Gateway Lan deletion to be executed
```

## Examples

```text
ionosctl natgateway lan remove --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID
```

