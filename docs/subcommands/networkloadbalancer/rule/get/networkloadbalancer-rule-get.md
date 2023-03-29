---
description: Get a Network Load Balancer Forwarding Rule
---

# NetworkloadbalancerRuleGet

## Usage

```text
ionosctl networkloadbalancer rule get [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `rule` command:

```text
[r forwardingrule]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id

## Options

```text
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
  -i, --rule-id string                  The unique ForwardingRule Id (required)
```

## Examples

```text
ionosctl nlb rule get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID
```

