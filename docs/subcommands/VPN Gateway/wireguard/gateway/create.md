---
description: "Create a WireGuard Gateway"
---

# VpnWireguardGatewayCreate

## Usage

```text
ionosctl vpn wireguard gateway create [flags]
```

## Aliases

For `wireguard` command:

```text
[wg]
```

For `gateway` command:

```text
[g gw]
```

For `create` command:

```text
[c post]
```

## Description

Create the IONOS side of a WireGuard VPN.

Networking. The gateway attaches to --lan-id inside --datacenter-id (both in --location). --gateway-ip is a public IPv4 from an IPBlock in that location — the address remote peers dial. --connection-ip is the gateway's own private address on the LAN (CIDR). --interface-ip is the address of the WireGuard tunnel interface itself.

Keys. WireGuard is key-based: supply the gateway's PRIVATE key with either --private-key (inline) or --private-key-path (read from a file) — exactly one is required. The matching public key is generated and returned; hand it to each peer so they can trust this gateway.

--port is the UDP port the gateway listens on (default 51820); peers must target it.

Once AVAILABLE, register remote devices with 'vpn wireguard peer create'.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ID Name PublicKey Description GatewayIP InterfaceIPv4 InterfaceIPv6 DatacenterId LanId ConnectionIPv4 ConnectionIPv6 InterfaceIP ListenPort Status]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-ip string      The gateway's own private address on the LAN, in CIDR notation (IPv4 or IPv6), e.g. 10.7.222.100/24 (required)
      --datacenter-id string      ID of the Virtual Data Center holding the LAN the gateway attaches to (required)
  -D, --depth int                 Level of detail for response objects (default 1)
      --description string        Description of the WireGuard Gateway
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
      --gateway-ip string         Public IPv4 from an IPBlock in the same location as the datacenter; this is the address remote peers connect to (required)
  -h, --help                      Print usage
      --interface-ip string       Address (with CIDR mask) of the WireGuard tunnel interface itself, IPv4 or IPv6 (required)
      --lan-id string             Numeric ID of the LAN the gateway attaches to; the private networks it will route into the tunnel live here (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string               Name of the WireGuard Gateway (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --port int                  UDP port the gateway listens on; peers must target it (default 51820, the WireGuard standard) (default 51820)
  -K, --private-key string        Gateway's WireGuard private key, inline (exactly one of this or --private-key-path is required). The public key is derived and returned for peers to trust
  -k, --private-key-path string   Path to a file holding the gateway's WireGuard private key (alternative to --private-key)
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vpn wireguard gateway create --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --connection-ip CONNECTION_IP --gateway-ip GATEWAY_IP --interface-ip INTERFACE_IP --private-key PRIVATE_KEY 
```

