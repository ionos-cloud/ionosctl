---
description: "Find a peer by ID"
---

# VpnWireguardPeerGet

## Usage

```text
ionosctl vpn wireguard peer get [flags]
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

For `get` command:

```text
[g]
```

## Description

Find a peer by ID

## Options

```text
  -u, --api-url string      Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -i, --peer-id string      The ID of the WireGuard Peer you want to delete (required)
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn wg g delete ...
```

