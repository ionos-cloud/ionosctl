---
description: "Remove a IPSec Tunnel"
---

# VpnIpsecTunnelDelete

## Usage

```text
ionosctl vpn ipsec tunnel delete [flags]
```

## Aliases

For `tunnel` command:

```text
[p]
```

For `delete` command:

```text
[d del rm]
```

## Description

Remove a IPSec Tunnel

## Options

```text
  -a, --all                 Delete all tunnels. Required or --tunnel-id
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
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -i, --tunnel-id string    The ID of the IPSec Tunnel you want to delete (required)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl vpn ipsec tunnel delete --gateway-id GATEWAY_ID --tunnel-id TUNNEL_ID 
```

