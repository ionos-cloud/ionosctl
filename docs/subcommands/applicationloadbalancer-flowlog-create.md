---
description: Create a Application Load Balancer FlowLog
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

Use this command to create a Application Load Balancer FlowLog in a specified Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                       Specifies the traffic Action pattern (default "ALL")
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -d, --direction string                    Specifies the traffic Direction pattern (default "INGRESS")
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -n, --name string                         The name for the FlowLog (default "Unnamed ALB Flow Log")
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
  -b, --s3bucket string                     S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog creation [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --bucket-name BUCKET_NAME
```

