---
description: Attach a CD-ROM to a Server
---

# ServerCdromAttach

## Usage

```text
ionosctl server cdrom attach [flags]
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

For `attach` command:

```text
[a]
```

## Description

Use this command to attach a CD-ROM to an existing Server.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Cdrom attachment [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for CD-ROM attachment to be executed
```

## Examples

```text
ionosctl server cdrom attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request
```

