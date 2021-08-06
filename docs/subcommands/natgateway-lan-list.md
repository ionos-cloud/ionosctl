---
description: List NAT Gateway Lans
---

# NatgatewayLanList

## Usage

```text
ionosctl natgateway lan list [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list NAT Gateway Lans from a specified NAT Gateway.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps] (default [NatGatewayLanId,GatewayIps])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
      --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl natgateway lan list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

