---
description: "List Target Groups Targets"
---

# TargetgroupTargetList

## Usage

```text
ionosctl targetgroup target list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Target Groups Targets.

## Options

```text
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
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
ionosctl targetgroup target list --targetgroup-id TARGET_GROUP_ID
```

