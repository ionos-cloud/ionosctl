---
description: "Update a Target Group"
---

# TargetgroupUpdate

## Usage

```text
ionosctl compute targetgroup update [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `update` command:

```text
[u up]
```

## Description

Update a Target Group's distribution, health-check, or HTTP health-check settings.

Only the flags you pass are changed; unspecified settings keep their current values. This command does NOT manage the targets (backend servers) themselves - use the `target` sub-commands (add/remove) for that.

Changes propagate to every ALB forwarding rule that references this group, so tightening the health check or switching --algorithm affects live traffic distribution.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Target Group Id

## Options

```text
      --algorithm string        How traffic is distributed across targets. ROUND_ROBIN: served alternately, honoring weights. LEAST_CONNECTION: the target with the fewest active connections is served next. RANDOM: chosen by a consistent pseudo-random function. SOURCE_IP: the same client IP always reaches the same target (source-based stickiness). (default "ROUND_ROBIN")
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --check-interval int      [Connection Health Check] Interval in milliseconds between consecutive health checks. Default is 2000. (default 2000)
      --check-timeout int       [Connection Health Check] Maximum time in milliseconds to wait for a target to respond to a check. If a target also has --check-interval set, the smaller of the two values is used once the TCP connection is established. (default 2000)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval State Retries Path Method MatchType Response Regex Negate]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Maximum number of items to return per request (default 50)
      --match-type string       [HTTP Health Check] What part of the target's reply decides health. STATUS_CODE: --response is matched against the HTTP status code. RESPONSE_BODY: --response is matched against the response body. (default "STATUS_CODE")
      --method string           [HTTP Health Check] The HTTP method used for the health check request. (default "GET")
  -n, --name string             The name of the target group. Used only for display; does not need to be unique. (default "Updated Target Group")
      --negate                  [HTTP Health Check] Invert the match: the target is healthy when --response does NOT match. Default is false.
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --path string             [HTTP Health Check] The request path (URL) the check sends to each target, e.g. /healthz. Default is '/'. (default "/.")
  -p, --protocol string         The forwarding protocol. Only HTTP is currently supported by the API. (default "HTTP")
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
      --regex                   [HTTP Health Check] Treat --response as a regular expression when matching the response body, instead of a literal value. Default is false.
      --response string         [HTTP Health Check] The value a target must return to be considered healthy. Interpreted per --match-type: a status code (e.g. 200) for STATUS_CODE, or expected body text for RESPONSE_BODY. (default "200")
      --retries int             [Connection Health Check] Maximum number of reconnection attempts to a target after a connection failure before it is marked unhealthy. Valid range is 0 to 65535; default is 3. (default 3)
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a group
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --name new-name -w

# Switch the algorithm to source-IP stickiness
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --algorithm SOURCE_IP

# Retune the HTTP health check to match a body instead of a status code
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --match-type RESPONSE_BODY --response OK
```

