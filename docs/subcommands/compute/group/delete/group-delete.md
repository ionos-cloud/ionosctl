---
description: Delete a Group
---

# GroupDelete

## Usage

```text
ionosctl group delete [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `delete` command:

```text
[d]
```

## Description

Use this operation to delete a single Group. Resources that are assigned to the Group are NOT deleted, but are no longer accessible to the Group members unless the member is a Contract Owner, Admin, or Resource Owner.

Required values to run command:

* Group Id

## Options

```text
  -a, --all                Delete all Groups.
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -i, --group-id string    The unique Group Id (required)
  -t, --timeout int        Timeout option for Request for Group deletion [seconds] (default 60)
  -w, --wait-for-request   Wait for Request for Group deletion to be executed
```

## Examples

```text
ionosctl group delete --group-id GROUP_ID
```

