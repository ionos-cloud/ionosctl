---
description: Add a Http Rule to Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleHttpruleAdd

## Usage

```text
ionosctl applicationloadbalancer rule httprule add [flags]
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

For `add` command:

```text
[a]
```

## Description

Use this command to add a Http Rule in a specified Application Load Balancer Forwarding Rule. 

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name
* Http Rule Type

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Name Type TargetGroupId DropQuery Location StatusCode ResponseMessage ContentType Condition] (default [Name,Type,TargetGroupId,DropQuery,Condition])
  -C, --condition string                    Matching rule for the Http Rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST and COOKIE types; must be null when type is SOURCE_IP (default "STARTS_WITH")
  -K, --condition-key string                Must be null when type is PATH, METHOD, HOST or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, QUERY (default "forward-at")
  -T, --condition-type string               Type of the Http Rule condition (default "HEADER")
  -V, --condition-value string              Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP (default "Friday")
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --content-type string                 Valid only for action STATIC (default "application/json")
      --datacenter-id string                The unique Data Center Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -l, --location string                     The location for redirecting; mandatory and valid only for REDIRECT action (default "www.ionos.com")
  -m, --message string                      The response message of the request; mandatory for STATIC action (default "Application Down")
  -n, --name string                         The name of a Application Load Balancer http rule; unique per forwarding rule (required)
      --negate                              Specifies whether the condition is negated or not; default: false
  -o, --output string                       Desired output format [text|json] (default "text")
  -Q, --query                               Default is false; must be true for REDIRECT action
  -q, --quiet                               Quiet output
      --rule-id string                      The unique ForwardingRule Id (required)
      --status-code int                     Valid only for action REDIRECT and STATIC; on REDIRECT action default is 301 and it can take the values 301, 302, 303, 307, 308; on STATIC action default is 503 and it can take a value between 200 and 599 (default 301)
      --targetgroup-id string               The ID of the target group; mandatory and only valid for FORWARD action
  -t, --timeout int                         Timeout option for Request for Forwarding Rule Http Rule creation [seconds] (default 300)
      --type string                         Type of the Http Rule (required)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule Http Rule creation to be executed
```

## Examples

```text
ionosctl alb rule http add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME --type TYPE
```
