---
description: "Create a IPSec tunnel"
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

Create IPSec tunnels

## Options

```text
  -u, --api-url string                    Override default host URL (default "https://vpn.de-fra.ionos.com")
      --auth-method string                The authentication method for the IPSec tunnel. Valid values are PSK or RSA (required)
      --cloud-network-cidrs strings       The network CIDRs on the "Left" side that are allowed to connect to the IPSec tunnel, i.e the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses.
      --cols strings                      Set of columns to be printed on output 
                                          Available columns: [ID Name Description RemoteHost AuthMethod PSKKey IKEDiffieHellmanGroup IKEEncryptionAlgorithm IKEIntegrityAlgorithm IKELifetime ESPDiffieHellmanGroup ESPEncryptionAlgorithm ESPIntegrityAlgorithm ESPLifetime CloudNetworkCIDRs PeerNetworkCIDRs Status StatusMessage]
  -c, --config string                     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --description string                Description of the IPSec Tunnel
      --esp-diffie-hellman-group string   The Diffie-Hellman Group to use for IPSec Encryption.. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --esp-encryption-algorithm string   The encryption algorithm to use for IPSec Encryption.. Can be one of: AES128-CTR, AES256-CTR, AES128-GCM-16, AES256-GCM-16, AES128-GCM-12, AES256-GCM-12, AES128-CCM-12, AES256-CCM-12, AES128, AES256
      --esp-integrity-algorithm string    The integrity algorithm to use for IPSec Encryption.. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --esp-lifetime int32                The phase lifetime in seconds
  -f, --force                             Force command to execute without user input
  -i, --gateway-id string                 The ID of the IPSec Gateway (required)
  -h, --help                              Print usage
      --host string                       The remote peer host fully qualified domain name or IPV4 IP to connect to. * __Note__: This should be the public IP of the remote peer. * Tunnels only support IPV4 or hostname (fully qualified DNS name). (required)
      --ike-diffie-hellman-group string   The Diffie-Hellman Group to use for IPSec Encryption.. Can be one of: 15-MODP3072, 16-MODP4096, 19-ECP256, 20-ECP384, 21-ECP521, 28-ECP256BP, 29-ECP384BP, 30-ECP512BP
      --ike-encryption-algorithm string   The encryption algorithm to use for IPSec Encryption.. Can be one of: AES128, AES256
      --ike-integrity-algorithm string    The integrity algorithm to use for IPSec Encryption.. Can be one of: SHA256, SHA384, SHA512, AES-XCBC
      --ike-lifetime int32                The phase lifetime in seconds
      --json-properties string            Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example           If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
  -l, --location string                   Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci
      --name string                       Name of the IPSec Tunnel (required)
      --no-headers                        Don't print table headers when table output is used
  -o, --output string                     Desired output format [text|json|api-json] (default "text")
      --peer-network-cidrs strings        The network CIDRs on the "Right" side that are allowed to connect to the IPSec tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses.
      --psk-key string                    The pre-shared key for the IPSec tunnel (required)
  -q, --quiet                             Quiet output
  -t, --timeout int                       Timeout in seconds for polling the request (default 60)
  -v, --verbose                           Print step-by-step process when running command
  -w, --wait                              Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl vpn ipsec tunnel create --gateway-id GATEWAY_ID --name NAME --host HOST --auth-method AUTH_METHOD --psk-key PSK_KEY --ike-diffie-hellman-group IKE_DIFFIE_HELLMAN_GROUP --ike-encryption-algorithm IKE_ENCRYPTION_ALGORITHM --ike-integrity-algorithm IKE_INTEGRITY_ALGORITHM --ike-lifetime IKE_LIFETIME --esp-diffie-hellman-group ESP_DIFFIE_HELLMAN_GROUP --esp-encryption-algorithm ESP_ENCRYPTION_ALGORITHM --esp-integrity-algorithm ESP_INTEGRITY_ALGORITHM --esp-lifetime ESP_LIFETIME --cloud-network-cidrs CLOUD_NETWORK_CIDRS --peer-network-cidrs PEER_NETWORK_CIDRS 
ionosctl vpn ipsec tunnel create --json-properties JSON_PROPERTIES 
ionosctl vpn ipsec tunnel create --json-properties JSON_PROPERTIES  json-properties-example
```

