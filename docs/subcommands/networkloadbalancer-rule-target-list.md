---
description: List Network Load Balancer Forwarding Rule Targets
---

# NetworkloadbalancerRuleTargetList

## Usage

```text
ionosctl networkloadbalancer rule target list [flags]
```

## Aliases

For `rule` command:
```text
[r forwarding]
```

For `target` command:
```text
[t]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to list Targets of a Network Load Balancer Forwarding Rule.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -u, --api-url string                  Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance] (default [TargetIp,TargetPort,Weight,Check,CheckInterval,Maintenance])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -f, --force                           Force command to execute without user input
  -h, --help                            help for list
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
      --rule-id string                  The unique ForwardingRule Id (required)
```

## Examples

```text
ionosctl networkloadbalancer rule target list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID
```

