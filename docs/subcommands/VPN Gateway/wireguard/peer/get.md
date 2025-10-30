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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the WireGuard Gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -i, --peer-id string      The ID of the WireGuard Peer (required)
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn wg peer get --gateway-id GATEWAY_ID --peer-id PEER_ID 
```

