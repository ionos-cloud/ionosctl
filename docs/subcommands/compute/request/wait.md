---
description: Wait a Request
---

# RequestWait

## Usage

```text
ionosctl request wait [flags]
```

## Aliases

For `request` command:

```text
[req]
```

For `wait` command:

```text
[w]
```

## Description

Use this command to wait for a specified Request to execute.

You can specify a timeout for the Request to be executed using `--timeout` option.

Required values to run command:

* Request Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -i, --request-id string   The unique Request Id (required)
  -t, --timeout int         Timeout option waiting for Request [seconds] (default 60)
```

## Examples

```text
ionosctl request wait --request-id REQUEST_ID
```

