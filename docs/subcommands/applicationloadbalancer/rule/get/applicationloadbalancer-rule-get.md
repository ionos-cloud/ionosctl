---
description: Get a Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleGet

## Usage

```text
ionosctl applicationloadbalancer rule get [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
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

Use this command to get information about a specified Application Load Balancer Forwarding Rule from a Application Load Balancer.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -i, --rule-id string                      The unique ForwardingRule Id (required)
```

## Examples

```text
ionosctl alb rule get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID
```

