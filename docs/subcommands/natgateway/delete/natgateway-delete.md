---
description: Delete a NAT Gateway
---

# NatgatewayDelete

## Usage

```text
ionosctl natgateway delete [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified NAT Gateway from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
  -a, --all                    Delete all Natgateways.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --natgateway-id string   The unique NatGateway Id (required)
  -t, --timeout int            Timeout option for Request for NAT Gateway deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway deletion to be executed
```

## Examples

```text
ionosctl natgateway delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```
