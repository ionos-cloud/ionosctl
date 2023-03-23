---
description: Delete an Application Load Balancer FlowLog
---

# ApplicationloadbalancerFlowlogDelete

## Usage

```text
ionosctl applicationloadbalancer flowlog delete [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
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

Use this command to delete a specified Application Load Balancer FlowLog from an Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id

## Options

```text
  -a, --all                                 Delete all Application Load Balancer FlowLogs
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string                   The unique FlowLog Id (required)
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog deletion [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog deletion to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog delete --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID
```

