---
description: Update a LAN
---

# Update

## Usage

```text
ionosctl lan update [flags]
```

## Description

Use this command to update a specified LAN.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for update
      --ignore-stdin           Force command to execute without user input
      --lan-id string          The unique LAN Id [Required flag]
      --lan-name string        The name of the LAN
      --lan-public             Public option for LAN
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --timeout int            Timeout option [seconds] (default 60)
      --wait                   Wait for LAN to be updated
```

## Examples

```text
ionosctl lan update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --lan-name demoLAN --lan-public=true
LanId   Name      Public
3       demoLAN   true
RequestId: 0a174dca-62b1-4360-aef8-89fd31c196f2
Status: Command lan update has been successfully executed
```

