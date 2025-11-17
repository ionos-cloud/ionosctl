---
description: "Remove a Target from a Target Group"
---

# TargetgroupTargetRemove

## Usage

```text
ionosctl targetgroup target remove [flags]
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

Use this command to delete the specified Target from Target Group.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port

## Options

```text
  -a, --all                     Delete all Target Group Targets
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
      --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --ip ip                   IP of a balanced target VM (required)
      --limit int               Maximum number of items to return per request (default 50)
      --no-headers              Don't print table headers when table output is used
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -P, --port int                Port of the balanced target service. (range: 1 to 65535) (required) (default 8080)
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group Target deletion [seconds] (default 60)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request        Wait for the Request for Target Group Target deletion to be executed
```

## Examples

```text
ionosctl targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT
```

