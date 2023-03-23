---
description: Get a FlowLog
---

# FlowlogGet

## Usage

```text
ionosctl flowlog get [flags]
```

## Aliases

For `flowlog` command:

```text
[fl]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information of a specified FlowLog from a NIC.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FlowLog Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --flowlog-id string      The unique FlowLog Id (required)
      --nic-id string          The unique NIC Id (required)
      --no-headers             When using text output, don't print headers
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl flowlog get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID
```

