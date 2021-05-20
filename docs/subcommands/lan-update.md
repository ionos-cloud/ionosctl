---
description: Update a LAN
---

# LanUpdate

## Usage

```text
ionosctl lan update [flags]
```

## Description

Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Private Cross-Connect.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output. Example: --format "ResourceId,Name" (default [LanId,Name,Public,PccId,State])
  -h, --help                   help for update
      --lan-id string          The unique LAN Id (required)
  -n, --name string            The name of the LAN
  -o, --output string          Desired output format [text|json] (default "text")
      --pcc-id string          The unique Id of the Private Cross-Connect the LAN will connect to
      --public                 Public option for LAN
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for LAN update to be executed
```

## Examples

```text
ionosctl lan update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --name demoLAN --public=true
LanId   Name      Public    PccId
3       demoLAN   true
RequestId: 0a174dca-62b1-4360-aef8-89fd31c196f2
Status: Command lan update has been successfully executed
```

