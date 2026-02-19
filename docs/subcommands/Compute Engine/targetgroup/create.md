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

Use this command to create a Target Group.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option.

## Options

```text
      --algorithm string     Balancing algorithm. (default "ROUND_ROBIN")
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --check-interval int   [Health Check] The interval in milliseconds between consecutive health checks; default is 2000. (default 2000)
      --check-timeout int    [Health Check] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established. (default 2000)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,CheckInterval,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --match-type string    [HTTP Health Check] Match Type for the HTTP health check. (default "STATUS_CODE")
      --method string        [HTTP Health Check] The method for the HTTP health check. (default "GET")
  -n, --name string          The name of the target group. (default "Unnamed Target Group")
      --negate               [HTTP Health Check] Negate for the HTTP health check.
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --path string          [HTTP Health Check] The path (destination URL) for the HTTP health check request; the default is /. (default "/.")
  -p, --protocol string      Balancing protocol (default "HTTP")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
      --regex                [HTTP Health Check] Regex for the HTTP health check.
      --response string      [HTTP Health Check] The response returned by the request, depending on the match type. (default "200")
      --retries int          [Health Check] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection attempts. (default 3)
  -t, --timeout int          Timeout option for Request for Target Group creation [seconds]. (default 60)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request     Wait for the Request for Target Group creation to be executed.
```

## Examples

```text
ionosctl compute targetgroup create --name TARGET_GROUP_NAME
```

