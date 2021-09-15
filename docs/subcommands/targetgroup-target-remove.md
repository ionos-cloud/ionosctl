---
description: Remove a Target from a Target Group
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
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance] (default [TargetIp,TargetPort,Weight,Check,CheckInterval,Maintenance])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -o, --output string           Desired output format [text|json] (default "text")
  -q, --quiet                   Quiet output
      --target-ip string        IP of a balanced target VM (required)
      --target-port int         Port of the balanced target service. (range: 1 to 65535) (required) (default 8080)
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group Target deletion [seconds] (default 60)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait-for-request        Wait for the Request for Target Group Target deletion to be executed
```

## Examples

```text
ionosctl targetgroup target remove --targetgroup-id TARGET_GROUP_ID --target-ip TARGET_IP --target-port TARGET_PORT
```

