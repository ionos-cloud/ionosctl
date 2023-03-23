---
description: Delete a Network Load Balancer FlowLog
---

# NetworkloadbalancerFlowlogDelete

## Usage

```text
ionosctl networkloadbalancer flowlog delete [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `flowlog` command:

```text
[f fl]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Network Load Balancer FlowLog from a Network Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id

## Options

```text
  -a, --all                             Delete all Network Load Balancer FlowLogs.
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string               The unique FlowLog Id (required)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -t, --timeout int                     Timeout option for Request for Network Load Balancer FlowLog deletion [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Network Load Balancer FlowLog deletion to be executed
```

## Examples

```text
ionosctl networkloadbalancer flowlog delete --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID
```

