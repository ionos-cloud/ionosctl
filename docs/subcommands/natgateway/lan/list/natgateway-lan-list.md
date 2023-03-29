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
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps] (default [NatGatewayLanId,GatewayIps])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl natgateway lan list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

