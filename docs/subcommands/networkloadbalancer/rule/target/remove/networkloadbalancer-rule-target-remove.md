---
description: Remove a Target from a Network Load Balancer Forwarding Rule
---

# NetworkloadbalancerRuleTargetRemove

## Usage

```text
ionosctl networkloadbalancer rule target remove [flags]
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

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a specified Target from Network Load Balancer Forwarding Rule.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port

## Options

```text
  -a, --all                             Remove all Forwarding Rule Targets.
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
      --ip ip                           IP of a balanced target VM (required)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -P, --port string                     Port of the balanced target service. Range: 1 to 65535 (required)
      --rule-id string                  The unique ForwardingRule Id (required)
  -t, --timeout int                     Timeout option for Request for Forwarding Rule Target deletion [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Forwarding Rule Target deletion to be executed
```

## Examples

```text
ionosctl nlb rule target remove --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip TARGET_IP --port TARGET_PORT
```

