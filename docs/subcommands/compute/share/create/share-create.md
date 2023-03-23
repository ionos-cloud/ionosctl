---
description: Create a Resource Share for a Group
---

# ShareCreate

## Usage

```text
ionosctl share create [flags]
```

## Aliases

For `create` command:

```text
[c]
```

## Description

Use this command to create a specific Resource Share to a Group and optionally allow the setting of permissions for that Resource. As an example, you might use this to grant permissions to use an Image or Snapshot to a specific Group.

Required values to run a command:

* Group Id
* Resource Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --edit-privilege       Set the group's permission to edit privileges on resource
      --group-id string      The unique Group Id (required)
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Set the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource to be shared through a Group [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Resource share to executed
```

## Examples

```text
ionosctl share create --group-id GROUP_ID --resource-id RESOURCE_ID
```

