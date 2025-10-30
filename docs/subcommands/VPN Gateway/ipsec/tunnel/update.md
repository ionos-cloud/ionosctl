---
description: "Update a IPSec Tunnel"
---

# VpnIpsecTunnelUpdate

## Usage

```text
ionosctl vpn ipsec tunnel update [flags]
```

## Aliases

For `tunnel` command:

```text
[p]
```

For `update` command:

```text
[u patch put]
```

## Description

Update a IPSec Tunnel

## Options

```text
  -u, --api-url string                    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'vpn' and env var 'IONOS_API_URL' (default "https://vpn.%s.ionos.com")
      --auth-method string                The authentication method for the IPSec tunnel. Valid values are PSK or RSA (required)
      --cloud-network-cidrs strings       The network CIDRs on the "Left" side that are allowed to connect to the IPSec tunnel, i.e the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses.
      --cols strings                      Set of columns to be printed on output 
                                          Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string                     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --description string                Description of the IPSec Tunnel
      --esp-diffie-hellman-group string   The Diffie-Hellman Group to use for IPSec Encryption.. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --esp-encryption-algorithm string   The encryption algorithm to use for IPSec Encryption.. Can be one of: AES128-CTR, AES256-CTR, AES128-GCM-16, AES256-GCM-16, AES128-GCM-12, AES256-GCM-12, AES128-CCM-12, AES256-CCM-12, AES128, AES256
      --esp-integrity-algorithm string    The integrity algorithm to use for IPSec Encryption.. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --esp-lifetime int32                The phase lifetime in seconds
  -f, --force                             Force command to execute without user input
      --gateway-id string                 The ID of the IPSec Gateway (required)
  -h, --help                              Print usage
      --host string                       The remote peer host fully qualified domain name or IPV4 IP to connect to. * __Note__: This should be the public IP of the remote peer. * Tunnels only support IPV4 or hostname (fully qualified DNS name). (required)
      --ike-diffie-hellman-group string   The Diffie-Hellman Group to use for IPSec Encryption.. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --ike-encryption-algorithm string   The encryption algorithm to use for IPSec Encryption.. Can be one of: AES128, AES256
      --ike-integrity-algorithm string    The integrity algorithm to use for IPSec Encryption.. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --ike-lifetime int32                The phase lifetime in seconds
      --json-properties string            Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example           If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
  -l, --location string                   Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
  -n, --name string                       Name of the IPSec Tunnel (required)
      --no-headers                        Don't print table headers when table output is used
  -o, --output string                     Desired output format [text|json|api-json] (default "text")
      --peer-network-cidrs strings        The network CIDRs on the "Right" side that are allowed to connect to the IPSec tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses.
      --psk-key string                    The pre-shared key for the IPSec tunnel (required)
  -q, --quiet                             Quiet output
  -i, --tunnel-id string                  The ID of the IPSec Tunnel (required)
  -v, --verbose count                     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vpn ipsec tunnel update --gateway-id GATEWAY_ID --tunnel-id TUNNEL_ID --name NAME 
```

