---
description: "List WireGuard Gateways"
---

# VpnWireguardGatewayList

## Usage

```text
ionosctl vpn wireguard gateway list [flags]
```

## Aliases

For `wireguard` command:

```text
[wg]
```

For `gateway` command:

```text
[g gw]
```

For `list` command:

```text
[l ls]
```

## Description

List WireGuard Gateways

## Options

```text
  -u, --api-url string      Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://vpn.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name PublicKey Description GatewayIP InterfaceIPv4 InterfaceIPv6 DatacenterId LanId ConnectionIPv4 ConnectionIPv6 InterfaceIP ListenPort Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn wireguard gateway list
```

