---
description: List NAT Gateways
---

# NatgatewayList

## Usage

```text
ionosctl natgateway list [flags]
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

Use this command to list NAT Gateways from a specified Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayId Name PublicIps State] (default [NatGatewayId,Name,PublicIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                see step by step process when running a command
```

## Examples

```text
ionosctl natgateway list --datacenter-id DATACENTER_ID
```

