---
description: Delete an IpBlock
---

# IpblockDelete

## Usage

```text
ionosctl ipblock delete [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -a, --all                 Delete all the IpBlocks.
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -t, --timeout int         Timeout option for Request for IpBlock deletion [seconds] (default 60)
  -w, --wait-for-request    Wait for the Request for IpBlock deletion to be executed
```

## Examples

```text
ionosctl ipblock delete --ipblock-id IPBLOCK_ID --wait-for-request
```

