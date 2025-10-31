---
description: "Update a Target Group"
---

# TargetgroupUpdate

## Usage

```text
ionosctl targetgroup update [flags]
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

Use this command to update a specified Target Group.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option.

Required values to run command:

* Target Group Id

## Options

```text
      --algorithm string        Balancing algorithm. (default "ROUND_ROBIN")
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --check-interval int      [Health Check] The interval in milliseconds between consecutive health checks; default is 2000. (default 2000)
      --check-timeout int       [Health Check] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established. (default 2000)
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,CheckInterval,State])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32             Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               Pagination limit: Maximum number of items to return per request (default 50)
      --match-type string       [HTTP Health Check] Match Type for the HTTP health check. (default "STATUS_CODE")
      --method string           [HTTP Health Check] The method for the HTTP health check. (default "GET")
  -n, --name string             The name of the target group. (default "Updated Target Group")
      --negate                  [HTTP Health Check] Negate for the HTTP health check.
      --no-headers              Don't print table headers when table output is used
      --offset int              Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --path string             [HTTP Health Check] The path (destination URL) for the HTTP health check request; the default is /. (default "/.")
  -p, --protocol string         Balancing protocol (default "HTTP")
  -q, --quiet                   Quiet output
      --regex                   [HTTP Health Check] Regex for the HTTP health check.
      --response string         [HTTP Health Check] The response returned by the request, depending on the match type. (default "200")
      --retries int             [Health Check] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection attempts. (default 3)
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group update [seconds]. (default 60)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request        Wait for the Request for Target Group update to be executed.
```

## Examples

```text
ionosctl targetgroup update --targetgroup-id TARGET_GROUP_ID --name TARGET_GROUP_NEW_NAME -w
```

