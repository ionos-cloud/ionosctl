---
description: "Remove a WireGuard Peer"
---

# VpnWireguardPeerDelete

## Usage

```text
ionosctl vpn wireguard peer delete [flags]
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

For `delete` command:

```text
[d del rm]
```

## Description

Remove a WireGuard Peer

## Options

```text
  -a, --all                 Delete all peers. Required or --peer-id
  -u, --api-url string      Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
      --location string     The location your resources are hosted in. Possible values: [de/fra de/txl] (default "de/txl")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -i, --peer-id string      The ID of the WireGuard Peer you want to delete (required)
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn wireguard peer delete --gateway-id GATEWAY_ID --peer-id PEER_ID 
```

