---
description: "Get a Target Group"
---

# TargetgroupGet

## Usage

```text
ionosctl targetgroup get [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Target Group.

Required values to run command:

* Target Group Id

## Options

```text
  -u, --api-url string          Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,CheckInterval,State])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32             Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -v, --verbose                 Print step-by-step process when running command
```

## Examples

```text
ionosctl targetgroup get -i TARGET_GROUP_ID
```

