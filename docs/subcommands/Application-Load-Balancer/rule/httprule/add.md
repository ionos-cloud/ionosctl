---
description: "Add an HTTP Rule to an Application Load Balancer Forwarding Rule"
---

# ApplicationloadbalancerRuleHttpruleAdd

## Usage

```text
ionosctl compute applicationloadbalancer rule httprule add [flags]
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

Use this command to add an HTTP rule to a forwarding rule (listener) on an Application Load Balancer. An HTTP rule matches incoming requests and then performs one action based on its --type:

  * FORWARD  - proxy matching requests to a backend target group (--targetgroup-id required).
  * REDIRECT - reply with an HTTP redirect to --location using --status-code (301/302/303/307/308); --query controls whether the original query string is dropped.
  * STATIC   - reply directly from the balancer with --status-code, --response-message and --content-type, without contacting any backend.

Matching is controlled by the condition flags: --condition-type selects which part of the request to inspect, --condition-key narrows it (e.g. a header name), --condition is the comparison operator, and --condition-value is what to compare against. Use --negate to invert the match. A rule with no conditions always matches and is useful as a default; rules are evaluated in order within the listener.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name
* Http Rule Type

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [Name Type TargetGroupId DropQuery Location StatusCode ResponseMessage ContentType Condition]
  -C, --condition string                    The comparison operator applied between the selected request element and --condition-value. Possible values: EXISTS, CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH. Mandatory for HEADER, PATH, QUERY, METHOD, HOST and COOKIE types; must be empty when condition-type is SOURCE_IP. (default "EQUALS")
  -K, --condition-key string                Narrows the condition to a specific named entry within the selected element, e.g. the header name "Accept" when condition-type=HEADER. Only valid for HEADER, COOKIE and QUERY; must be empty for PATH, METHOD, HOST and SOURCE_IP. (default "Accept")
  -T, --condition-type string               Selects which part of the incoming HTTP request the condition inspects. Possible values: HEADER, PATH, QUERY, METHOD, HOST, COOKIE, SOURCE_IP (default "HEADER")
  -V, --condition-value string              The value compared against the selected request element, e.g. "application/json" with condition=EQUALS, condition-type=HEADER, condition-key=Accept. Mandatory for CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH, and for condition-type SOURCE_IP (where it must be a valid CIDR); must be empty when condition is EXISTS. (default "application/json")
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --content-type string                 The Content-Type header of the static response, e.g. application/json or text/html. Only valid for --type STATIC. (default "application/json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Maximum number of items to return per request (default 50)
  -l, --location string                     The URL/host clients are redirected to. Mandatory and only valid for --type REDIRECT. (default "www.ionos.com")
  -m, --message string                      The response body the balancer returns. Mandatory and only used for --type STATIC. (default "Application Down")
  -n, --name string                         The unique name of the Application Load Balancer HTTP rule (also used to reference it when removing). (required)
      --negate                              Inverts the condition so the rule matches when the condition does NOT hold. Default is false.
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -Q, --query                               When true, drops the query string from the redirect target so the redirect URI carries no query parameters. Default is false; valid only for --type REDIRECT.
  -q, --quiet                               Quiet output
      --rule-id string                      The unique ForwardingRule Id (required)
      --status-code int                     The HTTP status code returned to the client. Only valid for REDIRECT and STATIC actions. For REDIRECT: 301, 302, 303, 307 or 308 (default 301). For STATIC: any value in the range 200-599 (API default 503). (default 301)
      --targetgroup-id string               The target group whose backend servers matching requests are proxied to. Mandatory and only valid for --type FORWARD.
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
      --type string                         The action the rule performs on matching requests: FORWARD (proxy to a target group), STATIC (reply directly from the balancer), or REDIRECT (send an HTTP redirect). (required)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# FORWARD every request to a target group (no conditions = catch-all)
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "to-backend" --type FORWARD --targetgroup-id TARGETGROUP_ID

# FORWARD only requests whose path starts with /api
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "api-route" --type FORWARD --targetgroup-id TARGETGROUP_ID --condition-type PATH --condition STARTS_WITH --condition-value /api

# REDIRECT all traffic to an HTTPS host with a 301
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "force-https" --type REDIRECT --location https://www.ionos.com --status-code 301

# STATIC maintenance page returned directly by the balancer
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "maintenance" --type STATIC --status-code 503 --content-type text/html --message "<h1>Under maintenance</h1>"
```

