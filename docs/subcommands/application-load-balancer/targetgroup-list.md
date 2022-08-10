---
description: List Target Groups
---

# TargetgroupList

## Usage

```text
ionosctl targetgroup list [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Target Groups.

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,CheckInterval,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings   Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -M, --max-results int   The maximum number of elements to return (default 500)
      --order-by string   Limits results to those containing a matching value for a specific property
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl targetgroup list
```

