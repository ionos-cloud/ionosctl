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

Create a WireGuard Gateway

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --connection-ip string   A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway (required)
      --datacenter-id string   The datacenter to connect your VPN Gateway to (required)
      --description string     Description of the WireGuard Gateway
  -f, --force                  Force command to execute without user input
      --gateway-ip string      Public IP address to be assigned to the gateway. Note: This must be an IP address in the same datacenter as the connections (required)
  -h, --help                   Print usage
      --interface-ip string    The IPv4 or IPv6 address (with CIDR mask) to be assigned to the WireGuard interface (required)
      --lan-id string          The numeric LAN ID to connect your VPN Gateway to (required)
  -n, --name string            Name of the WireGuard Gateway (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --port int               Port that WireGuard Server will listen on (default 51820)
  -k, --private-key string     The private key to be used by the WireGuard Gateway (required)
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl 
```

