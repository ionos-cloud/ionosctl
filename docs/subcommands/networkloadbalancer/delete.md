---
description: Delete a Network Load Balancer
---

# NetworkloadbalancerDelete

## Usage

```text
ionosctl networkloadbalancer delete [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
  -a, --all                             Delete all NetworkLoadBalancers.
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -i, --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -t, --timeout int                     Timeout option for Request for Network Load Balancer deletion [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Network Load Balancer deletion to be executed
```

## Examples

```text
ionosctl networkloadbalancer delete --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID
```

