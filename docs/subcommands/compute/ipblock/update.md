---
description: Update an IpBlock
---

# IpblockUpdate

## Usage

```text
ionosctl ipblock update [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -n, --name string         Name of the IpBlock
  -t, --timeout int         Timeout option for Request for IpBlock update [seconds] (default 60)
  -w, --wait-for-request    Wait for the Request for IpBlock update to be executed
```

## Examples

```text
ionosctl ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME
```

