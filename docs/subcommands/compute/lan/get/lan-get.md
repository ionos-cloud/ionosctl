---
description: Get a LAN
---

# LanGet

## Usage

```text
ionosctl lan get [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information of a given LAN.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --lan-id string          The unique LAN Id (required)
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl lan get --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

