---
description: Create an Application Load Balancer FlowLog
---

# ApplicationloadbalancerFlowlogCreate

## Usage

```text
ionosctl applicationloadbalancer flowlog create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create an Application Load Balancer FlowLog in a specified Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                       Specifies the traffic action pattern. (default "ALL")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string                    Specifies the traffic direction pattern. (default "INGRESS")
  -n, --name string                         The name of the Application Load Balancer FlowLog. (default "Unnamed ALB Flow Log")
  -b, --s3bucket string                     S3 bucket name of an existing IONOS Cloud S3 bucket. (required)
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog creation [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --bucket-name BUCKET_NAME
```

