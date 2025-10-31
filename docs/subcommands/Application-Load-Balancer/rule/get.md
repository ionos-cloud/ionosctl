---
description: "Get a Application Load Balancer Forwarding Rule"
---

# ApplicationloadbalancerRuleGet

## Usage

```text
ionosctl applicationloadbalancer rule get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Application Load Balancer Forwarding Rule from a Application Load Balancer.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           pagination limit: Maximum number of items to return per request (default 50)
      --no-headers                          Don't print table headers when table output is used
      --offset int                          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
  -i, --rule-id string                      The unique ForwardingRule Id (required)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl alb rule get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID
```

