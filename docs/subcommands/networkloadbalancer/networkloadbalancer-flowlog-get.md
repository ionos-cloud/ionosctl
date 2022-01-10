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
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -i, --flowlog-id string               The unique FlowLog Id (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
  -v, --verbose                         Print step-by-step process when running command
```

## Examples

```text
ionosctl networkloadbalancer flowlog get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID
```

