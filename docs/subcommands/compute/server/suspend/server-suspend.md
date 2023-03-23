---
description: Suspend a Cube Server
---

# ServerSuspend

## Usage

```text
ionosctl server suspend [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

## Description

Use this command to suspend a Cube Server. The operation can only be applied to Cube Servers. Note: The virtual machine will not be deleted.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server suspend [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server suspend to be executed
```

## Examples

```text
ionosctl server suspend --datacenter-id DATACENTER_ID -i SERVER_ID
```

