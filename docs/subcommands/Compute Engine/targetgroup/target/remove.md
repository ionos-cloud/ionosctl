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
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --ip ip                   IP of a balanced target VM (required)
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -P, --port int                Port of the balanced target service. (range: 1 to 65535) (required) (default 8080)
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group Target deletion [seconds] (default 60)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait-for-request        Wait for the Request for Target Group Target deletion to be executed
```

## Examples

```text
ionosctl targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT
```

