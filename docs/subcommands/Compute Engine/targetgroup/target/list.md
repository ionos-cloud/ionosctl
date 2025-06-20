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
  -u, --api-url string          Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
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
ionosctl targetgroup target list --targetgroup-id TARGET_GROUP_ID
```

