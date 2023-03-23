---
description: Create a LAN
---

# LanCreate

## Usage

```text
ionosctl lan create [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a new LAN within a Virtual Data Center on your account. The name, the public option and the Private Cross-Connect Id can be set.

NOTE: IP Failover is configured after LAN creation using an update command.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -n, --name string            The name of the LAN (default "Unnamed LAN")
      --pcc-id string          The unique Id of the Private Cross-Connect the LAN will connect to
  -p, --public                 Indicates if the LAN faces the public Internet (true) or not (false). E.g.: --public=true, --public=false
  -t, --timeout int            Timeout option for Request for LAN creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for LAN creation to be executed
```

## Examples

```text
ionosctl lan create --datacenter-id DATACENTER_ID --name NAME --public=true
```

