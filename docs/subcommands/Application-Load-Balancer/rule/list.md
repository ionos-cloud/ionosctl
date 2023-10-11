---
description: "List Application Load Balancer Forwarding Rules"
---

# ApplicationloadbalancerRuleList

## Usage

```text
ionosctl applicationloadbalancer rule list [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
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

Use this command to list Application Load Balancer Forwarding Rules from a specified Application Load Balancer.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings                     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -M, --max-results int32                   The maximum number of elements to return
      --order-by string                     Limits results to those containing a matching value for a specific property
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
  -v, --verbose                             Print step-by-step process when running command
```

## Examples

```text
ionosctl applicationloadbalancer rule list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID
```

