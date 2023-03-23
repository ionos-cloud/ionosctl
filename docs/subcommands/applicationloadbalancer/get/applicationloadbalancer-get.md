---
description: Get an Application Load Balancer
---

# ApplicationloadbalancerGet

## Usage

```text
ionosctl applicationloadbalancer get [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Application Load Balancer from a Virtual Data Center. You can also wait for Application Load Balancer to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int                         Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds] (default 300)
  -W, --wait-for-state                      Wait for specified Application Load Balancer to be in AVAILABLE state
```

## Examples

```text
ionosctl applicationloadbalancer get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID
```

