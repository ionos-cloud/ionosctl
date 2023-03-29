---
description: Force a hard reboot of a Server
---

# ServerReboot

## Usage

```text
ionosctl server reboot [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `reboot` command:

```text
[r]
```

## Description

Use this command to force a hard reboot of the Server. Do not use this method if you want to gracefully reboot the machine. This is the equivalent of powering off the machine and turning it back on.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server reboot [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server reboot to be executed
```

## Examples

```text
ionosctl server reboot --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

