---
description: "Remove a Http Rule from a Application Load Balancer Forwarding Rule"
---

# ApplicationloadbalancerRuleHttpruleRemove

## Usage

```text
ionosctl applicationloadbalancer rule httprule remove [flags]
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

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a specified Http Rule from Application Load Balancer Forwarding Rule.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name

## Options

```text
  -a, --all                                 Remove all HTTP Rules
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Name Type TargetGroupId DropQuery Location StatusCode ResponseMessage ContentType Condition] (default [Name,Type,TargetGroupId,DropQuery,Condition])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -n, --name string                         A name of that Application Load Balancer Http Rule (required)
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
      --rule-id string                      The unique ForwardingRule Id (required)
  -t, --timeout int                         Timeout option for Request for Forwarding Rule Http Rule deletion [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait                                Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl alb rule httprule remove --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME
```

