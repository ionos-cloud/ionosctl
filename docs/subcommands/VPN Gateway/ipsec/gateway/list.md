---
description: "List IPSec Gateways"
---

# VpnIpsecGatewayList

## Usage

```text
ionosctl vpn ipsec gateway list [flags]
```

## Aliases

For `gateway` command:

```text
[g gw]
```

For `list` command:

```text
[l ls]
```

## Description

List IPSec Gateways

## Options

```text
  -u, --api-url string      Override default host url (default "vpn.de-txl.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ID Name Description GatewayIP DatacenterId LanId ConnectionIPv4 ConnectionIPv6 Version Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --location string     The location your resources are hosted in. Possible values: [de/fra de/txl] (default "de/txl")
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl vpn ipsec gateway list
```

