---
description: Delete a Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleDelete

## Usage

```text
ionosctl applicationloadbalancer rule delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Application Load Balancer Forwarding Rule from a Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -a, --all                                 Delete all Forwarding Rules
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -i, --rule-id string                      The unique ForwardingRule Id (required)
  -t, --timeout int                         Timeout option for Request for Forwarding Rule deletion [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule deletion to be executed
```

## Examples

```text
ionosctl applicationloadbalancer rule delete --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID
```

