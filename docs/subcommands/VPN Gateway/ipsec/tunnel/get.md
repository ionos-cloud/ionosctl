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
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
      --gateway-id string   The ID of the IPSec Gateway (required)
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -i, --tunnel-id string    The ID of the IPSec Tunnel (required)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn ipsec tunnel get --gateway-id GATEWAY_ID --tunnel-id TUNNEL_ID 
```

