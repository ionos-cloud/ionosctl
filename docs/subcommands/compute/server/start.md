---
description: Start a Server
---

# ServerStart

## Usage

```text
ionosctl server start [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `start` command:

```text
[on]
```

## Description

Use this command to start a Server from a Virtual Data Center. If the Server's public IP was deallocated then a new IP will be assigned.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server start [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server start to be executed
```

## Examples

```text
ionosctl server start --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

