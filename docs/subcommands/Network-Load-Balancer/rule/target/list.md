---
description: "List Network Load Balancer Forwarding Rule Targets"
---

# NetworkloadbalancerRuleTargetList

## Usage

```text
ionosctl networkloadbalancer rule target list [flags]
```

## Aliases

For `rule` command:

```text
[r forwardingrule]
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
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance] (default [TargetIp,TargetPort,Weight,Check,CheckInterval,Maintenance])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Maximum number of items to return per request (default 50)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
      --rule-id string                  The unique ForwardingRule Id (required)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl networkloadbalancer rule target list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID
```

