---
description: Get a Network Load Balancer
---

# NetworkloadbalancerGet

## Usage

```text
ionosctl networkloadbalancer get [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Network Load Balancer from a Virtual Data Center. You can also wait for Network Load Balancer to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -i, --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
  -t, --timeout int                     Timeout option for waiting for Network Load Balancer to be in AVAILABLE state [seconds] (default 300)
  -W, --wait-for-state                  Wait for specified Network Load Balancer to be in AVAILABLE state
```

## Examples

```text
ionosctl networkloadbalancer get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID
```

