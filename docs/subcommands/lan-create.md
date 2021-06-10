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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
  -n, --name string            The name of the LAN (default "Unnamed LAN")
  -o, --output string          Desired output format [text|json] (default "text")
      --pcc-id string          The unique Id of the Private Cross-Connect the LAN will connect to
      --public                 Indicates if the LAN faces the public Internet (true) or not (false)
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for LAN creation to be executed
```

## Examples

```text
ionosctl lan create --datacenter-id DATACENTER_ID --name NAMEd
```

