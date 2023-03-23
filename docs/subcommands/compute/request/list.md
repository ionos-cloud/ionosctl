---
description: List Requests
---

# RequestList

## Usage

```text
ionosctl request list [flags]
```

## Aliases

For `request` command:

```text
[req]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list all Requests on your account.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [method headers body url]
* filter by metadata: [createdDate createdBy etag requestStatus]

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
      --latest int          Show latest N Requests. If it is not set, all Requests will be printed (deprecated)
  -M, --max-results int32   The maximum number of elements to return
      --method string       Show only the Requests with this method. E.g CREATE, UPDATE, DELETE (deprecated)
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl request list --latest N
```

