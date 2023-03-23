---
description: Delete a Resource Share from a Group
---

# ShareDelete

## Usage

```text
ionosctl share delete [flags]
```

## Aliases

For `delete` command:

```text
[d]
```

## Description

This command deletes a Resource Share from a specified Group.

Required values to run command:

* Resource Id
* Group Id

## Options

```text
  -a, --all                  Delete all the Resources Share from a specified Group.
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --group-id string      The unique Group Id (required)
  -i, --resource-id string   The unique Resource Id (required)
  -t, --timeout int          Timeout option for Request for Resource Share deletion [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Resource Share deletion to be executed
```

## Examples

```text
ionosctl share delete --group-id GROUP_ID --resource-id RESOURCE_ID --wait-for-request
```

