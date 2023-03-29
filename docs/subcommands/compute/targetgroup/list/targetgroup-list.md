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
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32   The maximum number of elements to return
      --order-by string     Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl targetgroup list
```

