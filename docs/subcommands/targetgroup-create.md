---
description: Create a Target Group
---

# TargetgroupCreate

## Usage

```text
ionosctl targetgroup create [flags]
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
      --algorithm string         Algorithm for the balancing (default "ROUND_ROBIN")
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --check-timeout int        [Health Check] It specifies the time (in milliseconds) for a target VM in this pool to answer the check. If a target VM has CheckInterval set and CheckTimeout is set too, then the smaller value of the two is used after the TCP connection is established (default 2000)
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout ConnectTimeout TargetTimeout Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,ConnectTimeout,TargetTimeout,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --connection-timeout int   [Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed (default 5000)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --match-type string        [HTTP Health Check] Match Type for the HTTP health check (default "STATUS_CODE")
      --method string            [HTTP Health Check] The method for the HTTP health check (default "GET")
  -n, --name string              Name of the Target Group (default "Unnamed Target Group")
      --negate                   [HTTP Health Check] Negate for the HTTP health check
  -o, --output string            Desired output format [text|json] (default "text")
      --path string              [HTTP Health Check] The path for the HTTP health check (default "/.")
  -p, --protocol string          Protocol of the balancing (default "HTTP")
  -q, --quiet                    Quiet output
      --regex                    [HTTP Health Check] Regex for the HTTP health check
      --response string          [HTTP Health Check] The response returned by the request (default "200")
      --retries int              [Health Check] The number of retries to perform on a target VM after a connection failure. (valid range: [0, 65535]) (default 3)
      --target-timeout int       [Health Check] The maximum inactivity time (in milliseconds) on the target VM side (default 50000)
  -t, --timeout int              Timeout option for Request for Target Group creation [seconds] (default 60)
  -v, --verbose                  Print step-by-step process when running command
  -w, --wait-for-request         Wait for the Request for Target Group creation to be executed
```

## Examples

```text
ionosctl targetgroup create --name TARGET_GROUP_NAME
```
