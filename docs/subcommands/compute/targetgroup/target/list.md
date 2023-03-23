---
description: List Target Groups Targets
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
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -i, --targetgroup-id string   The unique Target Group Id (required)
```

## Examples

```text
ionosctl targetgroup target list --targetgroup-id TARGET_GROUP_ID
```

