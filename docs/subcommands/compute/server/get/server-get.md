---
description: Get a Server
---

# ServerGet

## Usage

```text
ionosctl server get [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Server from a Virtual Data Center. You can also wait for Server to get in AVAILABLE state using `--wait-for-state` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for waiting for Server to be in AVAILABLE state [seconds] (default 60)
  -W, --wait-for-state         Wait for specified Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server get --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

