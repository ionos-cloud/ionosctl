---
description: Delete a NIC
---

# NicDelete

## Usage

```text
ionosctl nic delete [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a specified NIC.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
  -a, --all                    Delete all the Nics from a Server.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --nic-id string          The unique NIC Id (required)
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for NIC deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NIC deletion to be executed
```

## Examples

```text
ionosctl nic delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --force
```

