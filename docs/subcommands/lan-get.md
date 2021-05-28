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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -i, --lan-id string          The unique LAN Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl lan get --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

