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
* Name
* Direction
* Bucket Name

## Options

```text
  -a, --action string                   Specifies the traffic Action pattern (default "ALL")
  -u, --api-url string                  Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
  -b, --bucket-name string              S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -d, --direction string                Specifies the traffic Direction pattern (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            help for create
  -n, --name string                     The name for the FlowLog (required)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
  -t, --timeout int                     Timeout option for Request for Network Load Balancer FlowLog creation [seconds] (default 300)
  -w, --wait-for-request                Wait for the Request for Network Load Balancer FlowLog creation to be executed
```

## Examples

```text
ionosctl networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --bucket-name BUCKET_NAME
```

