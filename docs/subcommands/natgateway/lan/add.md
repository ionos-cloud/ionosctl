---
description: Add a NAT Gateway Lan
---

# NatgatewayLanAdd

## Usage

```text
ionosctl natgateway lan add [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `add` command:

```text
[a]
```

## Description

Use this command to add a NAT Gateway Lan in a specified NAT Gateway.

If IPs are not set manually, using `--ips` option, an IP will be automatically assigned. IPs must contain valid subnet mask. If user will not provide any IP then system will generate an IP with /24 subnet.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id

## Options

```text
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps] (default [NatGatewayLanId,GatewayIps])
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ips strings            Collection of Gateway IPs. If not set, it will automatically reserve public IPs
  -i, --lan-id int             The unique LAN Id (required) (default 1)
      --natgateway-id string   The unique NatGateway Id (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway Lan addition [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway Lan addition to be executed
```

## Examples

```text
ionosctl natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID

ionosctl natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID --ips IP_1,IP_2
```

