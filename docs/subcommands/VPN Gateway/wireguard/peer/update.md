---
description: "Update a WireGuard Peer"
---

# VpnWireguardPeerUpdate

## Usage

```text
ionosctl vpn wireguard peer update [flags]
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

For `update` command:

```text
[u patch put]
```

## Description

Update a WireGuard Peer

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ID Name Description Host Port WhitelistIPs PublicKey Status]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --description string   Description of the WireGuard Peer
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
      --gateway-id string    The ID of the WireGuard Gateway (required)
  -h, --help                 Print usage
      --host string          Hostname or IPV4 address that the WireGuard Server will connect to (required)
      --ips strings          Comma separated subnets of CIDRs that are allowed to connect to the WireGuard Gateway. Specify "a.b.c.d/32" for an individual IP address. Specify "0.0.0.0/0" or "::/0" for all addresses (required)
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
  -n, --name string          Name of the WireGuard Peer (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --peer-id string       The ID of the WireGuard Peer (required)
      --port int             Port that the WireGuard Server will connect to (default 51820)
      --public-key string    Public key of the connecting peer (required)
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn wireguard peer update --gateway-id GATEWAY_ID --peer-id PEER_ID --name NAME --description DESCRIPTION --ips IPS --public-key PUBLIC_KEY --host HOST --port PORT 
```

