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
  -a, --all                     Delete all Target Groups
  -D, --depth int32             Controls the detail depth of the response objects. Max depth is 10.
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group deletion [seconds] (default 60)
  -w, --wait-for-request        Wait for the Request for Target Group deletion to be executed
```

## Examples

```text
ionosctl targetgroup delete --targetgroup-id TARGET_GROUP_ID --force
```

