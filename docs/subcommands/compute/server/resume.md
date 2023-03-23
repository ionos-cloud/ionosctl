---
description: Resume a Cube Server
---

# ServerResume

## Usage

```text
ionosctl server resume [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `resume` command:

```text
[res]
```

## Description

Use this command to resume a Cube Server. The operation can only be applied to suspended Cube Servers.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server resume [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server resume to be executed
```

## Examples

```text
ionosctl server resume --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

