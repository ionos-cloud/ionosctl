---
description: "Update a WireGuard Gateway"
---

# VpnWireguardGatewayUpdate

## Usage

```text
ionosctl vpn wireguard gateway update [flags]
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

For `update` command:

```text
[u put patch]
```

## Description

Update a WireGuard Gateway

## Options

```text
  -u, --api-url string            Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ID Name Description GatewayIP InterfaceIPv4 InterfaceIPv6 DatacenterId LanId ConnectionIPv4 ConnectionIPv6 InterfaceIP ListenPort Status]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --connection-ip string      A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway (required)
      --datacenter-id string      The datacenter to connect your VPN Gateway to (required)
      --description string        Description of the WireGuard Gateway
  -f, --force                     Force command to execute without user input
  -i, --gateway-id string         The ID of the WireGuard Gateway (required)
      --gateway-ip string         The IP of an IPBlock in the same location as the provided datacenter (required)
  -h, --help                      Print usage
      --interface-ip string       The IPv4 or IPv6 address (with CIDR mask) to be assigned to the WireGuard interface (required)
      --lan-id string             The numeric LAN ID to connect your VPN Gateway to (required)
  -n, --name string               Name of the WireGuard Gateway (required)
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --port int                  Port that WireGuard Server will listen on (default 51820)
  -K, --private-key string        Specify the private key (required or --private-key-path)
  -k, --private-key-path string   Specify the private key from a file (required or --private-key)
  -q, --quiet                     Quiet output
  -v, --verbose                   Print step-by-step process when running command
```

