---
description: Update a Network Load Balancer FlowLog
---

# NetworkloadbalancerFlowlogUpdate

## Usage

```text
ionosctl networkloadbalancer flowlog update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Network Load Balancer FlowLog from a Network Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id

## Options

```text
  -a, --action string                   Specifies the traffic Action pattern
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string                Specifies the traffic Direction pattern
  -i, --flowlog-id string               The unique FlowLog Id (required)
  -n, --name string                     Name of the Network Load Balancer FlowLog
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -b, --s3bucket string                 S3 Bucket name of an existing IONOS Cloud S3 Bucket
  -t, --timeout int                     Timeout option for Request for Network Load Balancer FlowLog update [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Network Load Balancer FlowLog update to be executed
```

## Examples

```text
ionosctl networkloadbalancer flowlog update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID --name NAME
```

