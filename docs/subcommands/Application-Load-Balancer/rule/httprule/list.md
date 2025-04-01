---
description: "List Application Load Balancer Forwarding Rule Http Rules"
---

# ApplicationloadbalancerRuleHttpruleList

## Usage

```text
ionosctl applicationloadbalancer rule httprule list [flags]
```

## Aliases

For `rule` command:

```text
[r forwardingrule]
```

For `httprule` command:

```text
[http]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Http Rules of a Application Load Balancer Forwarding Rule.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Name Type TargetGroupId DropQuery Location StatusCode ResponseMessage ContentType Condition] (default [Name,Type,TargetGroupId,DropQuery,Condition])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -M, --max-results int32                   The maximum number of elements to return
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
      --rule-id string                      The unique ForwardingRule Id (required)
  -v, --verbose count                       Print step-by-step process when running command
```

## Examples

```text
ionosctl alb rule http list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID
```

