---
description: Create a Network Load Balancer FlowLog
---

# NetworkloadbalancerFlowlogCreate

## Usage

```text
ionosctl networkloadbalancer flowlog create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a Network Load Balancer FlowLog in a specified Network Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                   Specifies the traffic Action pattern (default "ALL")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string                Specifies the traffic Direction pattern (default "BIDIRECTIONAL")
  -n, --name string                     The name for the FlowLog (default "Unnamed FlowLog")
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -b, --s3bucket string                 S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
  -t, --timeout int                     Timeout option for Request for Network Load Balancer FlowLog creation [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Network Load Balancer FlowLog creation to be executed
```

## Examples

```text
ionosctl networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --s3bucket BUCKET_NAME
```

