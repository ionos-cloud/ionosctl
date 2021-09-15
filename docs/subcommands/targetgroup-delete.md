---
description: Delete a Target Group
---

# TargetgroupDelete

## Usage

```text
ionosctl targetgroup delete [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete the specified Target Group.

Required values to run command:

* Target Group Id

## Options

```text
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout ConnectTimeout TargetTimeout Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,ConnectTimeout,TargetTimeout,State])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
  -o, --output string           Desired output format [text|json] (default "text")
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group deletion [seconds] (default 60)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait-for-request        Wait for the Request for Target Group deletion to be executed
```

## Examples

```text
ionosctl targetgroup delete --targetgroup-id TARGET_GROUP_ID --force
```

