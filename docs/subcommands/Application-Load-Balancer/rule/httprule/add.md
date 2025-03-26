---
description: "Add a Http Rule to Application Load Balancer Forwarding Rule"
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
  -C, --condition string                    comparison rule for condition-value and the element selected with condition-type and condition-key. Possible values: EXISTS, CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH. mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP. (default "EQUALS")
  -K, --condition-key string                selects which entry in the selected HTTP element is used for the rule. For example, "Accept" for condition-type=HEADER. Not applicable for HOST, PATH or SOURCE_IP (default "Accept")
  -T, --condition-type string               selects which element in the incoming HTTP request is used for the rule. Possible values HEADER, PATH, QUERY, METHOD, HOST, COOKIE, SOURCE _IP (default "HEADER")
  -V, --condition-value string              value compared with the selected HTTP element. For example "application/json" in combination with condition=EQUALS, condition-type=HEADER, condition-key=Accept would be valid. Not applicable for condition EXISTS. Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS (default "application/json")
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --content-type string                 Valid only for STATIC actions. (default "application/json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -l, --location string                     The location for redirecting; mandatory and valid only for REDIRECT actions. (default "www.ionos.com")
  -m, --message string                      The response message of the request; mandatory for STATIC actions. (default "Application Down")
  -n, --name string                         The unique name of the Application Load Balancer HTTP rule. (required)
      --negate                              Specifies whether the condition is negated or not
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -Q, --query                               Default is false; valid only for REDIRECT actions.
  -q, --quiet                               Quiet output
      --rule-id string                      The unique ForwardingRule Id (required)
      --status-code int                     Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599. (default 301)
      --targetgroup-id string               he ID of the target group; mandatory and only valid for FORWARD actions.
  -t, --timeout int                         Timeout in seconds for polling the request (default 60)
      --type string                         Type of the HTTP rule. (required)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait                                Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl alb rule http add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME --type TYPE
```

