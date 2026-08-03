---
description: "Add a Target to a Target Group"
---

# TargetgroupTargetAdd

## Usage

```text
ionosctl compute targetgroup target add [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `target` command:

```text
[t]
```

For `add` command:

```text
[a]
```

## Description

Add a backend server (target) to a Target Group. A target is identified by its --ip and --port; the same IP with different ports counts as distinct targets.

--weight controls this target's share of traffic relative to the others. Each target receives load proportional to its weight over the sum of all weights (so weight 2 gets twice the traffic of weight 1). Range is 0-256, default 1. A weight of 0 excludes the target from new load-balancing decisions but still lets it serve existing persistent connections - useful for gracefully draining a server. When sizing by capacity, start in the middle of the range (e.g. 10-100) so you can adjust up or down later.

--health-check-enabled (default true) decides whether this target is probed at all. When off, the target is treated as always available and traffic is sent to it blindly. When on, the target only receives traffic while it passes the group's health check (a connection attempt to the target's own IP and port).

--maintenance-enabled (default false) takes the target out of rotation regardless of health, so no balanced traffic reaches it.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port

## Options

```text
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
      --health-check-enabled    When true (default), the target only receives traffic while it passes the group's health check (a TCP connection attempt to this target's IP and port). When false, the target is treated as always available and is never probed. (default true)
  -h, --help                    Print usage
      --ip ip                   The IP address of the backend server that will receive balanced traffic. (required)
      --limit int               Maximum number of items to return per request (default 50)
  -m, --maintenance-enabled     When true, the target is held out of rotation and receives no balanced traffic regardless of its health status. Default is false.
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -P, --port int                The port on the backend server that receives traffic. Valid range is 1 to 65535. Together with --ip it uniquely identifies the target. (required) (default 8080)
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
  -W, --weight int              This target's share of traffic relative to the combined weight of all targets (higher weight = larger share). Valid range is 0 to 256, default 1. Weight 0 excludes the target from new balancing decisions but still serves existing persistent connections (useful for draining). Prefer mid-range values (e.g. 10-100) to leave room for later adjustment. (default 1)
```

## Examples

```text
# Add a backend with default weight 1 and health checking on
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.5 --port 8080

# Add a higher-capacity backend that should take twice the traffic
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.6 --port 8080 --weight 2

# Add a backend already in maintenance (registered but not receiving traffic)
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.7 --port 8080 --maintenance-enabled
```

