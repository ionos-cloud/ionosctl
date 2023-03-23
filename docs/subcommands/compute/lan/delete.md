---
description: Delete a LAN
---

# LanDelete

## Usage

```text
ionosctl lan delete [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -a, --all                    Delete all Lans from a Virtual Data Center.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --lan-id string          The unique LAN Id (required)
  -t, --timeout int            Timeout option for Request for LAN deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for LAN deletion to be executed
```

## Examples

```text
ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID

ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID --wait-for-request
```

