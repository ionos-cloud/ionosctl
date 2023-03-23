---
description: Delete an Application Load Balancer
---

# ApplicationloadbalancerDelete

## Usage

```text
ionosctl applicationloadbalancer delete [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -a, --all                                 Delete all Application Load Balancers
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int                         Timeout option for Request for Application Load Balancer deletion [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer deletion to be executed
```

## Examples

```text
ionosctl applicationloadbalancer delete --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID
```

