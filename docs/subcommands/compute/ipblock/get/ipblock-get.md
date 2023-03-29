---
description: Get an IpBlock
---

# IpblockGet

## Usage

```text
ionosctl ipblock get [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the attributes of a specific IpBlock.

Required values to run command:

* IpBlock Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -i, --ipblock-id string   The unique IpBlock Id (required)
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl ipblock get --ipblock-id IPBLOCK_ID
```

