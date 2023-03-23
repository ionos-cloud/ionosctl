---
description: List Resources Shares through a Group
---

# ShareList

## Usage

```text
ionosctl share list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a full list of all the Resources that are shared through a specified Group.

Required values to run command:

* Group Id

## Options

```text
  -a, --all                 List all resources without the need of specifying parent ID name.
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
      --group-id string     The unique Group Id (required)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl share list --group-id GROUP_ID
```

