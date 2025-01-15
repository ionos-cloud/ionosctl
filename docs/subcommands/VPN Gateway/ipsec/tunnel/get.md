---
description: "Find a tunnel by ID"
---

# VpnIpsecTunnelGet

## Usage

```text
ionosctl vpn ipsec tunnel get [flags]
```

## Aliases

For `tunnel` command:

```text
[p]
```

For `get` command:

```text
[g]
```

## Description

Find a tunnel by ID

## Options

```text
  -u, --api-url string      Override default host URL (default "https://vpn.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the IPSec Gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -i, --tunnel-id string    The ID of the IPSec Tunnel (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn ipsec tunnel get --gateway-id GATEWAY_ID --tunnel-id TUNNEL_ID 
```
