---
description: Remove a Http Rule from a Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleHttpruleRemove

## Usage

```text
ionosctl applicationloadbalancer rule httprule remove [flags]
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

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a specified Http Rule from Application Load Balancer Forwarding Rule.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name

## Options

```text
  -a, --all                                 Remove all HTTP Rules
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -n, --name string                         A name of that Application Load Balancer Http Rule (required)
      --rule-id string                      The unique ForwardingRule Id (required)
  -t, --timeout int                         Timeout option for Request for Forwarding Rule Http Rule deletion [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule Http Rule deletion to be executed
```

## Examples

```text
ionosctl alb rule httprule remove --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME
```

