---
description: "Create a WireGuard Peer"
---

# VpnWireguardPeerCreate

## Usage

```text
ionosctl vpn wireguard peer create [flags]
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

For `create` command:

```text
[c post]
```

## Description

Add a remote device (peer) to a WireGuard gateway.

Supply the peer's --public-key so the gateway trusts it, and --ips: the source subnets the peer may send through the tunnel (WireGuard routes by key + allowed IPs). Use "a.b.c.d/32" for a single host, or "0.0.0.0/0","::/0" to allow everything.

--host/--port tell the gateway where to reach the peer for outbound connections; for peers behind NAT that only dial in, point --host at any reachable address (WireGuard learns the real endpoint from incoming traffic).

There is a per-gateway limit on peers; see product documentation.

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
  -i, --gateway-id string    The ID of the WireGuard Gateway (required)
  -h, --help                 Print usage
      --host string          Hostname or IPv4 the gateway uses to reach this peer (for peers behind NAT, any reachable address; the real endpoint is learned from inbound traffic) (required)
      --ips strings          Comma-separated CIDRs the peer is allowed to send through the tunnel (its allowed source IPs). Use "a.b.c.d/32" for a single host, or "0.0.0.0/0","::/0" for all addresses (required)
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string          Name of the WireGuard Peer (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --port int             UDP port the gateway uses to reach this peer (default 51820) (default 51820)
      --public-key string    The peer's WireGuard public key; the gateway trusts the device holding the matching private key (required)
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vpn wireguard peer create --gateway-id GATEWAY_ID --name my-laptop --public-key PUBLIC_KEY --ips 10.7.222.0/24 --host vpn.example.com
ionosctl vpn wireguard peer create --gateway-id GATEWAY_ID --name allow-all --public-key PUBLIC_KEY --ips 0.0.0.0/0,::/0 --host 203.0.113.5 --port 51820
```

