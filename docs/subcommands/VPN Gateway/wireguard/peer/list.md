---
description: "List WireGuard Peers"
---

# VpnWireguardPeerList

## Usage

```text
ionosctl vpn wireguard peer list [flags]
```

## Aliases

For `wireguard` command:

```text
[wg]
```

For `peer` command:

```text
[p]
```

For `list` command:

```text
[l ls]
```

## Description

List WireGuard Peers

## Options

```text
  -u, --api-url string      Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
      --location string     The location your resources are hosted in. Possible values: [de/fra de/txl] (default "de/txl")
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn wireguard peer list --gateway-id GATEWAY_ID 
```

