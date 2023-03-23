---
description: Delete a Network Load Balancer Forwarding Rule
---

# NetworkloadbalancerRuleDelete

## Usage

```text
ionosctl networkloadbalancer rule delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -a, --all                             Delete all Network Load Balancer Forwarding Rule.
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -i, --rule-id string                  The unique ForwardingRule Id (required)
  -t, --timeout int                     Timeout option for Request for Forwarding Rule deletion [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Forwarding Rule deletion to be executed
```

## Examples

```text
ionosctl networkloadbalancer rule delete --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID
```

