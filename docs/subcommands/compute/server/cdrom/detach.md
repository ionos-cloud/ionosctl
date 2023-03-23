---
description: Detach a CD-ROM from a Server
---

# ServerCdromDetach

## Usage

```text
ionosctl server cdrom detach [flags]
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

For `detach` command:

```text
[d]
```

## Description

This will detach the CD-ROM from the Server.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -a, --all                    Detach all CD-ROMS from a Server.
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for CD-ROM detachment [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for CD-ROM detachment to be executed
```

## Examples

```text
ionosctl server cdrom detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request --force
```

