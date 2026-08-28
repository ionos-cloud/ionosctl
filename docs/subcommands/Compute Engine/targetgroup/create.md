---
description: "Create a Target Group"
---

# TargetgroupCreate

## Usage

```text
ionosctl compute targetgroup create [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `create` command:

```text
[c]
```

## Description

Create a Target Group: a reusable backend pool that an Application Load Balancer forwarding rule can point to.

This command creates the group and its health-check configuration only. It starts with zero targets; add backend servers afterwards with `ionosctl compute targetgroup target add`. The group does nothing on its own until an ALB FORWARD rule references its ID.

The group defines two independent layers of health checking that must both pass for a target to receive traffic:
  - Connection check (--check-timeout, --check-interval, --retries): a TCP-level probe.
  - HTTP check (--path, --method, --match-type, --response): an application-level probe whose success is decided by matching either the HTTP status code or the response body.

--protocol only accepts HTTP. See --algorithm for how traffic is distributed across targets.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

## Options

```text
      --algorithm string     How traffic is distributed across targets. ROUND_ROBIN: served alternately, honoring weights. LEAST_CONNECTION: the target with the fewest active connections is served next. RANDOM: chosen by a consistent pseudo-random function. SOURCE_IP: the same client IP always reaches the same target (source-based stickiness). (default "ROUND_ROBIN")
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --check-interval int   [Connection Health Check] Interval in milliseconds between consecutive health checks. Default is 2000. (default 2000)
      --check-timeout int    [Connection Health Check] Maximum time in milliseconds to wait for a target to respond to a check. If a target also has --check-interval set, the smaller of the two values is used once the TCP connection is established. (default 2000)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval State Retries Path Method MatchType Response Regex Negate]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --match-type string    [HTTP Health Check] What part of the target's reply decides health. STATUS_CODE: --response is matched against the HTTP status code. RESPONSE_BODY: --response is matched against the response body. (default "STATUS_CODE")
      --method string        [HTTP Health Check] The HTTP method used for the health check request. (default "GET")
  -n, --name string          The name of the target group. Used only for display; does not need to be unique. (default "Unnamed Target Group")
      --negate               [HTTP Health Check] Invert the match: the target is healthy when --response does NOT match. Default is false.
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --path string          [HTTP Health Check] The request path (URL) the check sends to each target, e.g. /healthz. Default is '/.'. (default "/.")
  -p, --protocol string      The forwarding protocol. Only HTTP is currently supported by the API. (default "HTTP")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
      --regex                [HTTP Health Check] Treat --response as a regular expression when matching the response body, instead of a literal value. Default is false.
      --response string      [HTTP Health Check] The value a target must return to be considered healthy. Interpreted per --match-type: a status code (e.g. 200) for STATUS_CODE, or expected body text for RESPONSE_BODY. (default "200")
      --retries int          [Connection Health Check] Maximum number of reconnection attempts to a target after a connection failure before it is marked unhealthy. Valid range is 0 to 65535; default is 3. (default 3)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Minimal group (defaults: ROUND_ROBIN, HTTP, GET / expecting status 200)
ionosctl compute targetgroup create --name web-backends

# Round-robin group whose HTTP health check expects status 200 on GET /healthz, waiting until AVAILABLE
ionosctl compute targetgroup create --name web-backends --algorithm ROUND_ROBIN --path /healthz --method GET --match-type STATUS_CODE --response 200 -w

# Source-IP sticky group with a body match and a stricter, faster health check
ionosctl compute targetgroup create --name api-backends --algorithm SOURCE_IP \
  --check-timeout 1500 --check-interval 1000 --retries 2 \
  --path /status --match-type RESPONSE_BODY --response OK
```

