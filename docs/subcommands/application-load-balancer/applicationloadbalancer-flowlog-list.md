---
description: List Application Load Balancer FlowLogs
---

# ApplicationloadbalancerFlowlogList

## Usage

```text
ionosctl applicationloadbalancer flowlog list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Application Load Balancer FlowLogs from a specified Application Load Balancer.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -F, --filters strings                     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -M, --max-results int                     The maximum number of elements to return
      --order-by string                     Limits results to those containing a matching value for a specific property
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
  -v, --verbose                             Print step-by-step process when running command
```

## Examples

```text
ionosctl applicationloadbalancer flowlog list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID
```

