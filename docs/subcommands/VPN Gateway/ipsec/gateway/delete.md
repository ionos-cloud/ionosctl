---
description: "Delete a gateway"
---

# VpnIpsecGatewayDelete

## Usage

```text
ionosctl vpn ipsec gateway delete [flags]
```

## Aliases

For `gateway` command:

```text
[g gw]
```

For `delete` command:

```text
[del d]
```

## Description

Delete a gateway

## Options

```text
  -a, --all                 Delete all gateways. Required or --gateway-id
  -u, --api-url string      Override default host URL. If set, this will be preferred over the location flag as well as the config file override. If unset, the default will only be used as a fallback (default "https://vpn.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the IPSec Gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci (default "de/fra")
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn ipsec gateway --gateway-id GATEWAY_ID 
```

