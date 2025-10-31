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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the Wireguard Gateway (required)
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn wireguard peer list --gateway-id GATEWAY_ID 
```

