---
description: List IP Failovers groups from a LAN
---

# IpfailoverList

## Usage

```text
ionosctl ipfailover list [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of IP Failovers groups from a LAN.

Required values to run command:

* Data Center Id
* Lan Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
      --lan-id string          The unique LAN Id (required)
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

