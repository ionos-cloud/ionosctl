---
description: "Find a gateway by ID"
---

# VpnWireguardGatewayGet

## Usage

```text
ionosctl vpn wireguard gateway get [flags]
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

For `get` command:

```text
[g]
```

## Description

Find a gateway by ID

## Options

```text
  -u, --api-url string      Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name PublicKey Description GatewayIP InterfaceIPv4 InterfaceIPv6 DatacenterId LanId ConnectionIPv4 ConnectionIPv6 InterfaceIP ListenPort Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
      --location string     The location your resources are hosted in. Possible values: [de/fra de/txl] (default "de/txl")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn wireguard gateway get --gateway-id GATEWAY_ID 
```

