---
description: "Create a IPSec Gateway"
---

# VpnIpsecGatewayCreate

## Usage

```text
ionosctl vpn ipsec gateway create [flags]
```

## Aliases

For `gateway` command:

```text
[g gw]
```

For `create` command:

```text
[c post]
```

## Description

Create a IPSec Gateway

## Options

```text
  -u, --api-url string         Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-ip string   A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway (required)
      --datacenter-id string   The datacenter to connect your VPN Gateway to (required)
      --description string     Description of the IPSec Gateway
  -f, --force                  Force command to execute without user input
      --gateway-ip string      The IP of an IPBlock in the same location as the provided datacenter (required)
  -h, --help                   Print usage
      --lan-id string          The numeric LAN ID to connect your VPN Gateway to (required)
      --limit int              pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string        Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
  -n, --name string            Name of the IPSec Gateway (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --version string         The IKE version that is permitted for the VPN tunnels (default "IKEv2")
```

## Examples

```text
ionosctl vpn ipsec gateway create --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --connection-ip CONNECTION_IP --gateway-ip GATEWAY_IP --interface-ip INTERFACE_IP 
```

