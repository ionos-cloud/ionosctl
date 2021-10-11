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
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [NetworkLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -i, --flowlog-id string               The unique FlowLog Id (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
  -t, --timeout int                     Timeout option for Request for Network Load Balancer FlowLog deletion [seconds] (default 300)
  -v, --verbose                         Print step-by-step process when running command
  -w, --wait-for-request                Wait for the Request for Network Load Balancer FlowLog deletion to be executed
```

## Examples

```text
ionosctl networkloadbalancer flowlog delete --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID
```

