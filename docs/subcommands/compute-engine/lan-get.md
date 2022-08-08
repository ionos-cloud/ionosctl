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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --lan-id string          The unique LAN Id (required)
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl lan get --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

