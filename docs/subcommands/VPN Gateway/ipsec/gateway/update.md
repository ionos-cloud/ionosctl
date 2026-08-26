---
description: "Update an IPSec Gateway"
---

# VpnIpsecGatewayUpdate

## Usage

```text
ionosctl vpn ipsec gateway update [flags]
```

## Aliases

For `gateway` command:

```text
[g gw]
```

For `update` command:

```text
[u put patch]
```

## Description

Update an IPSec Gateway. This is a full replace (PUT): the gateway is fetched, your flags are applied on top, and the result is written back. Crypto and remote-host settings live on the tunnels, not here; see field meanings under 'ipsec gateway create'.

## Options

```text
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-ip string   The gateway's own private address on the LAN, in CIDR notation (IPv4 or IPv6), e.g. 10.7.222.100/24 (required)
      --datacenter-id string   ID of the Virtual Data Center holding the LAN the gateway attaches to (required)
  -D, --depth int              Level of detail for response objects (default 1)
      --description string     Description of the IPSec Gateway
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -i, --gateway-id string      The ID of the IPSec Gateway (required)
      --gateway-ip string      Public IPv4 from an IPBlock in the same location as the datacenter; this is the address remote peers connect to (required)
  -h, --help                   Print usage
      --lan-id string          Numeric ID of the LAN the gateway attaches to; the private networks it will route into the tunnel live here (required)
      --limit int              Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string            Name of the IPSec Gateway (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --version string         IKE version permitted for the tunnels on this gateway (currently only IKEv2) (default "IKEv2")
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vpn ipsec gateway update --gateway-id GATEWAY_ID --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --connection-ip CONNECTION_IP --gateway-ip GATEWAY_IP 
```

