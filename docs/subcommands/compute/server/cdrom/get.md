---
description: Get a specific attached CD-ROM from a Server
---

# ServerCdromGet

## Usage

```text
ionosctl server cdrom get [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `cdrom` command:

```text
[cd]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about an attached CD-ROM on Server.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server cdrom get --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID
```

