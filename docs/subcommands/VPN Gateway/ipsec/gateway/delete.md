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
  -u, --api-url string      Override default host URL (default "https://vpn.de-fra.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -i, --gateway-id string   The ID of the IPSec Gateway (required)
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, fr/par, gb/lhr, gb/bhx, us/ewr, us/las, us/mci
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout duration    Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl vpn ipsec gateway --gateway-id GATEWAY_ID 
```

