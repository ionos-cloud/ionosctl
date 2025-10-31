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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the IPSec Gateway (required)
  -h, --help                Print usage
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -i, --tunnel-id string    The ID of the IPSec Tunnel you want to delete (required)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn ipsec tunnel delete --gateway-id GATEWAY_ID --tunnel-id TUNNEL_ID 
```

