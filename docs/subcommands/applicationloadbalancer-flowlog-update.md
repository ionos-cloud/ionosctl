---
description: Update a Application Load Balancer FlowLog
---

# ApplicationloadbalancerFlowlogUpdate

## Usage

```text
ionosctl applicationloadbalancer flowlog update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Application Load Balancer FlowLog from a Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id

## Options

```text
  -a, --action string                       Specifies the traffic Action pattern
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -d, --direction string                    Specifies the traffic Direction pattern
  -i, --flowlog-id string                   The unique FlowLog Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                help for update
  -n, --name string                         Name of the Application Load Balancer FlowLog
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
  -b, --s3bucket string                     S3 Bucket name of an existing IONOS Cloud S3 Bucket
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog update [seconds] (default 300)
  -v, --verbose                             see step by step process when running a command
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog update to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID --name NAME
```

