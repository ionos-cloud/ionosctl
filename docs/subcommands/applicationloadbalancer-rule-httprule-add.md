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
                                            Available columns: [Name Type TargetGroupId DropQuery Location StatusCode ResponseMessage ContentType Condition] (default [Name,Type,TargetGroupId,DropQuery,Location,StatusCode,ResponseMessage,ContentType,Condition])
      --condition string                    Condition of the Http Rule condition
      --condition-key string                 (default "forward-at")
      --condition-type string               Type of the Http Rule condition (default "HEADER")
      --condition-value string               (default "Friday")
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --content-type string                  (default "application/json")
      --datacenter-id string                The unique Data Center Id (required)
      --drop-query                          Default is false; must be true for REDIRECT action
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -l, --location string                     The location for redirecting; mandatory for REDIRECT action (default "www.ionos.com")
  -n, --name string                         A name of that Application Load Balancer Http Rule (required) (default "Unnamed Http Rule")
      --negate                              Specifies whether the condition is negated or not; default: false
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
      --response string                     The response message of the request; mandatory for STATIC action
      --rule-id string                      The unique ForwardingRule Id (required)
      --status-code string                  On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599 (default "301")
      --targetgroup-id string               The Id of the Target Group; mandatory for FORWARD action
  -t, --timeout int                         Timeout option for Request for Forwarding Rule Http Rule creation [seconds] (default 300)
      --type string                         Type of the Http Rule (required)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule Http Rule creation to be executed
```

