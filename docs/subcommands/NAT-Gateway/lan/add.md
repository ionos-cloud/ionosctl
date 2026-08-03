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

Use this command to attach a private LAN to a NAT Gateway so servers on that LAN can route their outbound traffic through it. The gateway becomes reachable on the LAN via the gateway IPs given in `--ips` (the next-hop address servers use to reach the internet).

If `--ips` is not set, a gateway IP is generated automatically (with a /24 subnet). Gateway IPs must include a valid subnet mask and should belong to the same subnet as the LAN.

Attaching the LAN does not by itself translate traffic; add SNAT rules (`natgateway rule create`) whose `--source-subnet` covers the servers on this LAN.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

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
      --ips strings            Comma-separated gateway IPs (with subnet mask, e.g. 10.0.1.1/24) that the gateway uses on this LAN as the servers' next hop. Should belong to the LAN's subnet. If omitted, an IP is auto-generated with a /24 subnet
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
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Attach LAN 1 and let the gateway IP be auto-assigned
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id 1

# Attach LAN 1 with explicit gateway IPs (include the subnet mask)
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id 1 --ips 10.0.1.1/24
```

