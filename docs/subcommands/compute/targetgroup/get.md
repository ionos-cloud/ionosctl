---
description: Get a Target Group
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
  -D, --depth int32             Controls the detail depth of the response objects. Max depth is 10.
  -i, --targetgroup-id string   The unique Target Group Id (required)
```

## Examples

```text
ionosctl targetgroup get -i TARGET_GROUP_ID
```

