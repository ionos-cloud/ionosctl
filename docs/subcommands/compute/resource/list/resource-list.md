---
description: List Resources
---

# ResourceList

## Usage

```text
ionosctl resource list [flags]
```

## Aliases

For `resource` command:

```text
[res]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl resource list
```

