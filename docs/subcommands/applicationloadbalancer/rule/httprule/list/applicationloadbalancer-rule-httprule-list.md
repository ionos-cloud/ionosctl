---
description: List Application Load Balancer Forwarding Rule Http Rules
---

# ApplicationloadbalancerRuleHttpruleList

## Usage

```text
ionosctl applicationloadbalancer rule httprule list [flags]
```

## Aliases

For `rule` command:

```text
[r forwardingrule]
```

For `httprule` command:

```text
[http]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Http Rules of a Application Load Balancer Forwarding Rule.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -M, --max-results int32                   The maximum number of elements to return
      --rule-id string                      The unique ForwardingRule Id (required)
```

## Examples

```text
ionosctl alb rule http list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID
```

