---
description: List Network Load Balancer Forwarding Rules
---

# NetworkloadbalancerRuleList

## Usage

```text
ionosctl networkloadbalancer rule list [flags]
```

## Aliases

For `networkloadbalancer` command:
```text
[nlb]
```

For `rule` command:
```text
[r forwardingrule]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to list Network Load Balancer Forwarding Rules from a specified Network Load Balancer.

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
  -u, --api-url string                  Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [ForwardingRuleId Name Algorithm Protocol ListenerIp ListenerPort State ClientTimeout ConnectTimeout TargetTimeout Retries] (default [ForwardingRuleId,Name,Algorithm,Protocol,ListenerIp,ListenerPort,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            help for list
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
```

## Examples

```text
ionosctl networkloadbalancer rule list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID
```

