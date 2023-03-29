---
description: Create/Reserve an IpBlock
---

# IpblockCreate

## Usage

```text
ionosctl ipblock create [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

## Options

```text
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -l, --location string    Location of the IpBlock (default "de/txl")
  -n, --name string        Name of the IpBlock. If not set, it will automatically be set
      --size int           Size of the IpBlock (default 2)
  -t, --timeout int        Timeout option for Request for IpBlock creation [seconds] (default 60)
  -w, --wait-for-request   Wait for the Request for IpBlock creation to be executed
```

## Examples

```text
ionosctl ipblock create --name NAME --location LOCATION_ID --size IPBLOCK_SIZE
```

