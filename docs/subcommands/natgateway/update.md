---
description: Update a NAT Gateway
---

# NatgatewayUpdate

## Usage

```text
ionosctl natgateway update [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified NAT Gateway from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ips strings            Collection of public reserved IP addresses of the NAT Gateway. This will overwrite the current values
  -n, --name string            Name of the NAT Gateway
  -i, --natgateway-id string   The unique NatGateway Id (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway update to be executed
```

## Examples

```text
ionosctl natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME
```

