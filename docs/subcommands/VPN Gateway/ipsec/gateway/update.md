---
description: "Update a IPSec Gateway"
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

Update a IPSec Gateway

## Options

```text
  -u, --api-url string         Override default host URL (default "https://vpn.de-fra.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --connection-ip string   A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway (required)
      --datacenter-id string   The datacenter to connect your VPN Gateway to (required)
      --description string     Description of the IPSec Gateway
  -f, --force                  Force command to execute without user input
  -i, --gateway-id string      The ID of the IPSec Gateway (required)
      --gateway-ip string      The IP of an IPBlock in the same location as the provided datacenter (required)
  -h, --help                   Print usage
      --lan-id string          The numeric LAN ID to connect your VPN Gateway to (required)
  -l, --location string        Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci
  -n, --name string            Name of the IPSec Gateway (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for polling the request (default 60)
  -v, --verbose                Print step-by-step process when running command
      --version string         The IKE version that is permitted for the VPN tunnels (default "IKEv2")
  -w, --wait                   Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl vpn ipsec gateway update --gateway-id GATEWAY_ID --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --connection-ip CONNECTION_IP --gateway-ip GATEWAY_IP --interface-ip INTERFACE_IP 
```

