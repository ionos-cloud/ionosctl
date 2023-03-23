---
description: List Network Load Balancer Forwarding Rule Targets
---

# NetworkloadbalancerRuleTargetList

## Usage

```text
ionosctl networkloadbalancer rule target list [flags]
```

## Aliases

For `rule` command:

```text
[r forwardingrule]
```

For `target` command:

```text
[t]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Targets of a Network Load Balancer Forwarding Rule.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id

## Options

```text
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -M, --max-results int32               The maximum number of elements to return
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
      --rule-id string                  The unique ForwardingRule Id (required)
```

## Examples

```text
ionosctl networkloadbalancer rule target list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID
```

