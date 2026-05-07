---
description: "Add a NAT Gateway Lan"
---

# NatgatewayLanAdd

## Usage

```text
ionosctl compute natgateway lan add [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `add` command:

```text
[a]
```

## Description

Use this command to add a NAT Gateway Lan in a specified NAT Gateway.

If IPs are not set manually, using `--ips` option, an IP will be automatically assigned. IPs must contain valid subnet mask. If user will not provide any IP then system will generate an IP with /24 subnet.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            Collection of Gateway IPs. If not set, it will automatically reserve public IPs
  -i, --lan-id int             The unique LAN Id (required) (default 1)
      --limit int              Maximum number of items to return per request (default 50)
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes
```

## Examples

```text
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID --ips IP_1,IP_2
```

