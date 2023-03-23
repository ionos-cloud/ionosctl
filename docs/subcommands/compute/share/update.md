---
description: Update a Resource Share from a Group
---

# ShareUpdate

## Usage

```text
ionosctl share update [flags]
```

## Aliases

For `update` command:

```text
[u up]
```

## Description

Use this command to update the permissions that a Group has for a specific Resource Share.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --edit-privilege       Update the group's permission to edit privileges on resource
      --group-id string      The unique Group Id (required)
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Update the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource Share update [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Resource Share update to be executed
```

## Examples

```text
ionosctl share update --group-id GROUP_ID --resource-id RESOURCE_ID --share-privilege
```

