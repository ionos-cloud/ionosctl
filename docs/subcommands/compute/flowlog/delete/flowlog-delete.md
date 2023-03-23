---
description: Delete a FlowLog from a NIC
---

# FlowlogDelete

## Usage

```text
ionosctl flowlog delete [flags]
```

## Aliases

For `flowlog` command:

```text
[fl]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified FlowLog from a NIC.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FlowLog Id

## Options

```text
  -a, --all                    Delete all Flowlogs.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string      The unique FlowLog Id (required)
      --nic-id string          The unique NIC Id (required)
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for FlowLog deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for FlowLog deletion to be executed
```

## Examples

```text
ionosctl flowlog delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID -f -w
```

