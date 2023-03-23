---
description: Get a NAT Gateway
---

# NatgatewayGet

## Usage

```text
ionosctl natgateway get [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NAT Gateway from a Virtual Data Center. You can also wait for NAT Gateway to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             When using text output, don't print headers
  -t, --timeout int            Timeout option for waiting for NAT Gateway to be in AVAILABLE state [seconds] (default 60)
  -W, --wait-for-state         Wait for specified NAT Gateway to be in AVAILABLE state
```

## Examples

```text
ionosctl natgateway get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

