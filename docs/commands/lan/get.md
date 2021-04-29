---
description: Get a LAN
---

# Get

## Usage

```text
ionosctl lan get [flags]
```

## Description

Use this command to retrieve information of a specified LAN.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public,PccId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for get
      --lan-id string          The unique LAN Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl lan get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 4
LanId   Name      Public    PccId
4       demoLan   false
```

