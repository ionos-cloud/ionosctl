---
description: Get a Network Load Balancer FlowLog
---

# NetworkloadbalancerFlowlogGet

## Usage

```text
ionosctl networkloadbalancer flowlog get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Network Load Balancer FlowLog from a Network Load Balancer.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id

## Options

```text
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string               The unique FlowLog Id (required)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      When using text output, don't print headers
```

## Examples

```text
ionosctl networkloadbalancer flowlog get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID
```

