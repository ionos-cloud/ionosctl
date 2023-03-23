---
description: List IpConsumers
---

# IpconsumerList

## Usage

```text
ionosctl ipconsumer list [flags]
```

## Aliases

For `ipconsumer` command:

```text
[ipc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Resources where each IP address from an IpBlock is being used.

Required values to run command:

* IpBlock Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
      --ipblock-id string   The unique IpBlock Id (required)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl ipconsumer list --ipblock-id IPBLOCK_ID
```

