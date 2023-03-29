---
description: Get a Request
---

# RequestGet

## Usage

```text
ionosctl request get [flags]
```

## Aliases

For `request` command:

```text
[req]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Request.

Required values to run command:

* Request Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
      --no-headers          When using text output, don't print headers
  -i, --request-id string   The unique Request Id (required)
```

## Examples

```text
ionosctl request get --request-id REQUEST_ID
```

