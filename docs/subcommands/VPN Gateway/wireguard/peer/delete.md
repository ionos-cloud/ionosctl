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
  -u, --api-url string      Override default host URL (default "https://vpn.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -i, --peer-id string      The ID of the WireGuard Peer you want to delete (required)
  -q, --quiet               Quiet output
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl vpn wireguard peer delete --gateway-id GATEWAY_ID --peer-id PEER_ID 
```

