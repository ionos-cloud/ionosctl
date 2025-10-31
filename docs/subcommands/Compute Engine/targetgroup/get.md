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
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetGroupId Name Algorithm Protocol CheckTimeout CheckInterval Retries Path Method MatchType Response Regex Negate State] (default [TargetGroupId,Name,Algorithm,Protocol,CheckTimeout,CheckInterval,State])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32             Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --limit int               pagination limit: Maximum number of items to return per request (default 50)
      --no-headers              Don't print table headers when table output is used
      --offset int              pagination offset: Number of items to skip before starting to collect the results
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl targetgroup get -i TARGET_GROUP_ID
```

