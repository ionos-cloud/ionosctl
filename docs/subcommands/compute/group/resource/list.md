---
description: List Resources from a Group
---

# GroupResourceList

## Usage

```text
ionosctl group resource list [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `resource` command:

```text
[res]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Resources assigned to a Group. To see more details about existing Resources, use `ionosctl resource` commands.

Required values to run command:

* Group Id

## Options

```text
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ResourceId Name SecAuthProtection Type State] (default [ResourceId,Name,SecAuthProtection,Type,State])
      --group-id string     The unique Group Id (required)
  -M, --max-results int32   The maximum number of elements to return
```

## Examples

```text
ionosctl group resource list --group-id GROUP_ID
```

