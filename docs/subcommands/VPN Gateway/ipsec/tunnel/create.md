---
description: "Create an IPSec tunnel"
---

# VpnIpsecTunnelCreate

## Usage

```text
ionosctl vpn ipsec tunnel create [flags]
```

## Aliases

For `tunnel` command:

```text
[p]
```

For `create` command:

```text
[c post]
```

## Description

Create an IPSec tunnel: the connection from a gateway (--gateway-id) to one remote site.

Point --host at the remote peer's public IPv4 or FQDN, then authenticate: --auth-method PSK together with a shared --psk-key (or RSA). Set the phase-1 (IKE) and phase-2 (ESP) crypto with the --ike-* / --esp-* flags — each takes a Diffie-Hellman group, an encryption algorithm, an integrity algorithm and a lifetime in seconds (rekey interval; leave 0 to use the API default). Finally list which subnets may cross: --cloud-network-cidrs on your IONOS LAN side and --peer-network-cidrs on the remote side.

Both ends must use the SAME crypto parameters and mirrored CIDRs or the tunnel stays down.

You can instead pass the whole request body with --json-properties (see --json-properties-example for a template).

## Options

```text
  -u, --api-url string                    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --auth-method string                How the two ends authenticate each other: PSK (shared secret in --psk-key) or RSA (required)
      --cloud-network-cidrs strings       Local IONOS-side subnets (CIDR) allowed to cross the tunnel, i.e. the networks in your IONOS CLOUD LAN. Use "0.0.0.0/0","::/0" for all addresses
      --cols strings                      Set of columns to be printed on output 
                                          Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string                     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                         Level of detail for response objects (default 1)
      --description string                Description of the IPSec Tunnel
      --esp-diffie-hellman-group string   ESP (phase 2) Diffie-Hellman group for the data channel; must match the remote peer. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --esp-encryption-algorithm string   ESP (phase 2) encryption algorithm for the data channel; must match the remote peer. Can be one of: AES128-CTR, AES256-CTR, AES128-GCM-16, AES256-GCM-16, AES128-GCM-12, AES256-GCM-12, AES128-CCM-12, AES256-CCM-12, AES128, AES256
      --esp-integrity-algorithm string    ESP (phase 2) integrity/hash algorithm; must match the remote peer. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --esp-lifetime int32                ESP (phase 2) rekey interval in seconds; 0 uses the API default (e.g. 3600 = 1h)
  -F, --filters strings                   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                             Force command to execute without user input
  -i, --gateway-id string                 The ID of the IPSec Gateway (required)
  -h, --help                              Print usage
      --host string                       Public IPv4 or fully-qualified hostname of the remote peer to connect to (the remote side's public address; IPv6 is not supported) (required)
      --ike-diffie-hellman-group string   IKE (phase 1) Diffie-Hellman group for the key exchange; must match the remote peer. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --ike-encryption-algorithm string   IKE (phase 1) encryption algorithm; must match the remote peer. Can be one of: AES128, AES256
      --ike-integrity-algorithm string    IKE (phase 1) integrity/hash algorithm; must match the remote peer. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --ike-lifetime int32                IKE (phase 1) rekey interval in seconds; 0 uses the API default (e.g. 86400 = 24h)
      --json-properties string            Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example           If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                         Maximum number of items to return per request (default 50)
  -l, --location string                   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string                       Name of the IPSec Tunnel (required)
      --no-headers                        Don't print table headers when table output is used
      --offset int                        Number of items to skip before starting to collect the results
      --order-by string                   Property to order the results by
  -o, --output string                     Desired output format [text|json|api-json] (default "text")
      --peer-network-cidrs strings        Remote-side subnets (CIDR) reachable through the tunnel, i.e. the networks behind the remote peer. Use "0.0.0.0/0","::/0" for all addresses
      --psk-key string                    Pre-shared key, when --auth-method is PSK; the identical secret must be configured on the remote peer (required)
      --query string                      JMESPath query string to filter the output
  -q, --quiet                             Quiet output
  -t, --timeout int                       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vpn ipsec tunnel create --gateway-id GATEWAY_ID --name to-hq --host vpn.example.com --auth-method PSK --psk-key SHARED_SECRET --ike-diffie-hellman-group 16-MODP4096 --ike-encryption-algorithm AES256 --ike-integrity-algorithm SHA256 --esp-diffie-hellman-group 16-MODP4096 --esp-encryption-algorithm AES256 --esp-integrity-algorithm SHA256 --cloud-network-cidrs 10.7.222.0/24 --peer-network-cidrs 192.168.1.0/24
ionosctl vpn ipsec tunnel create --gateway-id GATEWAY_ID --json-properties tunnel.json
ionosctl vpn ipsec tunnel create --json-properties-example
```

