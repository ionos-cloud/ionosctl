---
description: Get a NIC
---

# NicGet

## Usage

```text
ionosctl nic get [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NIC from specified Data Center and Server.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --nic-id string          The unique NIC Id (required)
      --no-headers             When using text output, don't print headers
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl nic get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

