---
description: "Remove a Target from a Target Group"
---

# TargetgroupTargetRemove

## Usage

```text
ionosctl compute targetgroup target remove [flags]
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

For `remove` command:

```text
[r]
```

## Description

Remove a target (backend server) from a Target Group. The target is matched by its --ip and --port pair; both must match an existing target or the command reports it was not found. Use --all to remove every target from the group, leaving the group defined but empty.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port

## Options

```text
  -a, --all                     Remove all targets from the group, leaving it empty. Cannot be combined with --ip / --port.
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled]
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --ip ip                   The IP address of the target to remove. Must match an existing target's IP together with --port. (required)
      --limit int               Maximum number of items to return per request (default 50)
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -P, --port int                The port of the target to remove. Together with --ip it identifies which target is removed. Valid range is 1 to 65535. (required) (default 8080)
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                    Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Remove one backend by IP + port
ionosctl compute targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.5 --port 8080

# Empty the group (remove all targets)
ionosctl compute targetgroup target remove --targetgroup-id TARGET_GROUP_ID --all --force
```

